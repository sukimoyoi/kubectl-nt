# kubectl-neo-tree

kubectl treeを痒いところに手が届くようにした感じ

## Purpose

[kubectl tree](https://github.com/ahmetb/kubectl-tree)はリソースの親子関係を表示してくれるが、たとえばPVの表示はサポートしてなかったりする。kubectl-neo-tree(kubectl-nt)はkubectl treeがサポートしてないようなリソースの親子関係を表示する。


# Usgae

## install

```
go build cmd/kubectl-nt.go
cp ./kubectl-nt /usr/local/bin/
```


# How to use

```
$ kubectl nt pv pvc-251d1e11-e8ef-4d91-8bb6-9afeda2bf1ac
persistentvolume/pvc-251d1e11-e8ef-4d91-8bb6-9afeda2bf1ac (Parents)
└── persistentvolumeclaim/mybol-tssk-volume-prometheus-0

persistentvolume/pvc-251d1e11-e8ef-4d91-8bb6-9afeda2bf1ac (Children)
└── storageclass/ontap-block

$ kubectl nt sc ontap-block
storageclass/ontap-block (Parents)
└── persistentvolume/pvc-251d1e11-e8ef-4d91-8bb6-9afeda2bf1ac
│   ├── persistentvolumeclaim/mybol-tssk-volume-prometheus-0
└── persistentvolume/pvc-2b112065-0c4b-4374-9f62-701891f213db
│   ├── persistentvolumeclaim/mybol-tkks-volume-prometheus-0
└── persistentvolume/pvc-4ef8a8b1-dc90-4a0a-9c2d-a3fb5488c10b
│   ├── persistentvolumeclaim/mybol-pkks-volume-prometheus-0
└── persistentvolume/pvc-7c0625a0-4ae0-4d26-b1f7-1593efbcec36
│   ├── persistentvolumeclaim/mybol-pssk-volume-prometheus-0
└── persistentvolume/pvc-a5718168-3115-4426-8f1d-785af62be5d6
│   ├── persistentvolumeclaim/mybol-pssk-volume-prometheus-0
└── persistentvolume/pvc-abff4462-9fa8-46db-a672-fe38dd348c30
    └── persistentvolumeclaim/mybol-dssk-volume-prometheus-0
```

## support resources

plz see [relationer.go](./pkg/resourcerelationer/relationer.go#L10)

