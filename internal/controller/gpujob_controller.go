package controller

import (
    "context"
    "time"
    "log"

    gpuv1 "github.com/yourstartup/gpu-operator/api/v1"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/controller"
    "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type GpuJobReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop
func (r *GpuJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    var job gpuv1.GpuJob
    if err := r.Get(ctx, req.NamespacedName, &job); err != nil {
        log.Printf("unable to fetch GpuJob: %v", err)
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    start := time.Now()

    // Example: mark job as completed
    job.Status.Phase = "Completed"
    if err := r.Status().Update(ctx, &job); err != nil {
        log.Printf("failed to update status: %v", err)
        return ctrl.Result{}, err
    }

    // Record metrics
    GpuJobCounter.WithLabelValues(job.Namespace, string(job.Status.Phase)).Inc()
    GpuJobDuration.WithLabelValues(job.Namespace).Observe(time.Since(start).Seconds())

    return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GpuJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&gpuv1.GpuJob{}).
        WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
        Complete(r)
}

