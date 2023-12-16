## App Installer 

![CI](https://github.com/Drosaca/appImageInstaller/actions/workflows/go.yml/badge.svg)

<p >
  <a href="#">
    <img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white" alt="build"/>
  </a>
  <a href="https://github.com/Drosaca/appImageInstaller/releases">
    <img src="https://img.shields.io/github/v/release/Drosaca/appImageInstaller?style=for-the-badge&sort=semver" alt="tag"/>
  </a>
  <a href="https://github.com/d4data-official/d4data-app/issues">
    <img src="https://img.shields.io/github/issues-raw/Drosaca/appImageInstaller?style=for-the-badge" alt="issues"/>
  </a>
</p>


### Have you ever download an **AppImage** file ?

This type of binary has the advantage of being portable but it will not be present on your Desktop manager's applications list .

This tool written in Go allows to create app shrotcuts on the linux desktop environment and helps to better integrate appimages binaries to your system 

## Usage

```shell
sudo appinstall AppImageFilePath
```
After the installation, the Appimage is no longer needed and can be deleted.

## Demo
![](https://github.com/Drosaca/appImageInstaller/blob/main/assets/demo.gif)

## Install

Just download the binary file and copy it to your binaries path

```shell
sudo cp appinstall /usr/bin
```

That's it

_note: if you faced a bug feel free to open an issue_
