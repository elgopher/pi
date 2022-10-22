# How to install dependencies on Linux

## Debian/Ubuntu

```sh
sudo apt install gcc libc6-dev libglu1-mesa-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config
```

## Arch

```sh
sudo pacman -S gcc mesa libxrandr libxcursor libxinerama libxi pkg-config
```

## Fedora

```sh
sudo dnf install gcc mesa-libGLU-devel mesa-libGLES-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel libXxf86vm-devel alsa-lib-devel pkg-config
```
