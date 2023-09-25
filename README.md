# kpad

Kpad is a simple multiplatform terminal editor for editing kubernetes declarative manifest yaml files.
It provides a handy auto-complete that pops up a list of possible context-aware fields at the cursor position.

Behind the scenes, it calls "kubectl explain" to populate the auto-complete list, so it is also aware of your custom kubernetes objects in your cluster, and provides autocompletion for them.

If in your cluster you use something different from the plain "kubectl" CLI command, you can configure if in kpad launching `kpad -c` and change the configuration there.

For example if you use "microkubernetes", launch `kpad -c` and change like this:
```
kubectl: microk8s kubectl
```


Kpad is still work-in-progress, although it is pretty stable and quite fast.
It is written in go. 
It's pretty lightweight.
It compiles for Mac, Win, Linux.
It even runs on my phone.
It also supports wide characters.

There could be bugs, so you are welcome to contribute with fixes,enhancements,features,etc.


To make kpad your default kubernetes manifests editor, set the "KUBE_EDITOR" environment variable to the path where kpad is.
In linux, for example:
`export KUBE_EDITOR=/<path>/<to>/kpad`

