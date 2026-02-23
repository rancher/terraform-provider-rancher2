#!/bin/bash

kubectl get nodes || true
kubectl get all -A || true

kubectl get pods -A || true
sleep 10
kubectl get pods -A || true
sleep 10
kubectl get pods -A || true
