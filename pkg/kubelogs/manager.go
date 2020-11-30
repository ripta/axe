package kubelogs

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/ripta/axe/pkg/logger"
)

var ErrInformerNeverSynced = errors.New("informer cache never completed syncing")

type Manager struct {
	kubernetes.Interface

	l     logger.Interface
	logCh chan logger.LogLine
	mu    sync.Mutex

	nsCancelers     map[string]context.CancelFunc
	nsInformers     map[string]informers.SharedInformerFactory
	podLogCancelers map[string]context.CancelFunc

	containerTails map[string]bool

	lookback time.Duration
	resync   time.Duration
}

func NewManager(l logger.Interface, cs kubernetes.Interface, lookback, resync time.Duration) *Manager {
	if lookback > 0 {
		lookback = -lookback
	}
	if lookback == 0 {
		lookback = -5 * time.Minute
	}

	return &Manager{
		Interface: cs,
		l:         l,
		mu:        sync.Mutex{},
		logCh:     make(chan logger.LogLine, 1000),

		containerTails:  make(map[string]bool),
		nsCancelers:     make(map[string]context.CancelFunc),
		nsInformers:     make(map[string]informers.SharedInformerFactory),
		podLogCancelers: make(map[string]context.CancelFunc),

		lookback: lookback,
		resync:   resync,
	}
}

func (m *Manager) ContainerCount() (int, int) {
	var active, all int
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, up := range m.containerTails {
		all += 1
		if up {
			active += 1
		}
	}
	return active, all
}

func (m *Manager) Logs() <-chan logger.LogLine {
	return m.logCh
}

func (m *Manager) NamespaceCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.nsInformers)
}

func (m *Manager) Run(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for ns, inf := range m.nsInformers {
		ctx, cancel := context.WithCancel(ctx)

		m.nsCancelers[ns] = cancel
		inf.Start(ctx.Done())
	}

	return m.unsafeWaitForCacheSync(ctx.Done())
}

func (m *Manager) WaitForCacheSync(stopCh <-chan struct{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.unsafeWaitForCacheSync(stopCh)
}

func (m *Manager) unsafeWaitForCacheSync(stopCh <-chan struct{}) error {
	for ns, inf := range m.nsInformers {
		for typ, ok := range inf.WaitForCacheSync(stopCh) {
			if !ok {
				return fmt.Errorf("%w for type %s in namespace %s", ErrInformerNeverSynced, typ.String(), ns)
			}
		}

		m.l.Printf("cache synced for namespace %s", ns)
	}
	return nil
}

func (m *Manager) Watch(namespace string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.nsInformers[namespace]; !ok {
		inf := informers.NewSharedInformerFactoryWithOptions(m.Interface, m.resync, informers.WithNamespace(namespace))
		inf.Core().V1().Pods().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(newobj interface{}) {
				om, err := meta.Accessor(newobj)
				if err != nil {
					m.l.Printf("could not retrieve meta information from new object during add: %+v", err)
					return
				}
				m.startPodLogs(om.GetNamespace(), om.GetName())
			},
			UpdateFunc: func(_, newobj interface{}) {
				om, err := meta.Accessor(newobj)
				if err != nil {
					m.l.Printf("could not retrieve meta information from new object during update: %+v", err)
					return
				}
				m.startPodLogs(om.GetNamespace(), om.GetName())
			},
			DeleteFunc: func(oldobj interface{}) {
				om, err := meta.Accessor(oldobj)
				if err != nil {
					m.l.Printf("could not retrieve meta information from old object during delete: %+v", err)
					return
				}
				m.stopPodLogs(om.GetNamespace(), om.GetName())
			},
		})

		m.l.Printf("registered watch for namespace %s", namespace)
		m.nsInformers[namespace] = inf
	}
}

func (m *Manager) Unwatch(namespace string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.nsInformers[namespace]; ok {
		if cancel, ok := m.nsCancelers[namespace]; ok {
			cancel()
		}

		m.l.Printf("stopped watching namespace %s", namespace)
		delete(m.nsInformers, namespace)
		delete(m.nsCancelers, namespace)

		stops := make([]string, 0)
		for key := range m.podLogCancelers {
			if strings.HasPrefix(key, namespace+"/") {
				stops = append(stops, key)
			}
		}

		for _, key := range stops {
			if cancel, ok := m.podLogCancelers[key]; ok {
				cancel()
			}
			m.l.Printf("stopped tailing logs for %s", key)
			delete(m.podLogCancelers, key)
		}
	}
}

func (m *Manager) stopPodLogs(ns, name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.l.Printf("stopping pod logs for %s/%s", ns, name)

	key := fmt.Sprintf("%s/%s", ns, name)
	if cancel, ok := m.podLogCancelers[key]; ok {
		cancel()
	}
	delete(m.podLogCancelers, key)
}

func (m *Manager) startPodLogs(ns, name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%s/%s", ns, name)
	if _, ok := m.podLogCancelers[key]; ok {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.podLogCancelers[key] = cancel

	go m.tailPodLogs(ctx, ns, name)
}

func (m *Manager) tailPodLogs(ctx context.Context, ns, name string) {
	m.l.Printf("starting tail of logs for pod %s/%s", ns, name)
	defer m.stopPodLogs(ns, name)

	inf, ok := m.nsInformers[ns]
	if !ok {
		m.l.Printf("could not tail logs for %s/%s, because its namespace informer is missing", ns, name)
		return
	}

	pl := inf.Core().V1().Pods().Lister()
	pod, err := pl.Pods(ns).Get(name)
	if err != nil {
		if apierrors.IsNotFound(err) {
			m.l.Printf("ignoring deleted pod %s/%s", ns, name)
			return
		}
	}

	wg := sync.WaitGroup{}

	// TODO(ripta): handle init containers, which means we need to re-enter tailPodLogs
	for _, container := range pod.Spec.Containers {
		wg.Add(1)
		go m.tailPodContainerLogs(ctx, pl, ns, name, container.Name)
	}

	wg.Wait()
}

func (m *Manager) tailPodContainerLogs(ctx context.Context, pl listerv1.PodLister, ns, name, cn string) {
	key := fmt.Sprintf("%s/%s/%s", ns, name, cn)

	m.mu.Lock()
	m.containerTails[key] = true
	m.mu.Unlock()
	defer func() {
		m.mu.Lock()
		m.containerTails[key] = false
		m.mu.Unlock()
	}()

	m.l.Printf("starting tail of logs for container %s/%s/%s", ns, name, cn)
	plo := v1.PodLogOptions{
		Container: cn,
		Follow:    true,
		SinceTime: &metav1.Time{
			Time: time.Now().Add(m.lookback),
		},
	}

	for {
		if _, err := pl.Pods(ns).Get(name); err != nil {
			if apierrors.IsTooManyRequests(err) {
				// TODO(ripta): add jitter
				time.Sleep(5 * time.Second)
				m.l.Printf("got throttled by apiserver while asking about container %s/%s/%s", ns, name, cn)
				continue
			}
			if apierrors.IsNotFound(err) {
				m.l.Printf("ignoring container %s belonging to deleted pod %s/%s", cn, ns, name)
				return
			}
		}

		req := m.Interface.CoreV1().Pods(ns).GetLogs(name, &plo)
		stream, err := req.Context(ctx).Stream()
		if err != nil {
			// TODO(ripta): add jitter
			time.Sleep(5 * time.Second)
			m.l.Printf("could not tail %s/%s/%s: %+v", ns, name, cn, err)
			continue
		}
		defer stream.Close()

		// lag := time.NewTimer(time.Millisecond)
		// defer lag.Stop()

		m.l.Printf("streaming logs for container %s/%s/%s", ns, name, cn)
		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			line := logger.LogLine{
				Type:      logger.LogLineTypeContainer,
				Namespace: ns,
				Name:      name,
				Bytes:     scanner.Bytes(),
			}
			select {
			case m.logCh <- line:
			// case <-lag.C:
			// 	m.l.Printf("event buffer full, dropping logs for %s/%s", ns, name)
			case <-ctx.Done():
				m.l.Printf("stopped tailing %s/%s/%s", ns, name, cn)
				return
			}
		}

		plo.SinceTime.Time = time.Now()
		m.l.Printf("end of tail for container %s/%s/%s", ns, name, cn)
		if err := scanner.Err(); err != nil {
			m.l.Printf("error tailing container %s/%s/%s: %+v", ns, name, cn, err)
		}

		// TODO(ripta): add jitter
		time.Sleep(5 * time.Second)
	}
}
