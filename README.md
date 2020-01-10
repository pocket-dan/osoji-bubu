おそうじぶーぶー💨
--

〜 ラジコン視点で楽しくお掃除する 〜



おそうじぶーぶーはラジコンに乗せたカメラ映像をVRゴーグルで見ることで，あたかも自分がラジコンを運転しているような感覚を体験できるエンターテイメントです．さらに，ラジコンに掃除機を搭載することで”おそうじ”を可能にし，集めたゴミ量がスコアとしても可視化されるという作品です．



おそうじぶーぶーは，**「小さくなって冒険したい」という誰もが一度は抱く願望を叶え**ながら，**「おかたづけ，掃除が面倒くさい」という課題を解決** します！



<a href="https://youtu.be/HNjXZwRTybU?t=7517" target="_blank"><img src="https://user-images.githubusercontent.com/13511520/71416628-98e72a80-26a4-11ea-96a7-cf4255e14be1.png"></a>



## About

このレポジトリは掃除機（ラジコン）側に取り付けられたラズベリーパイに接続されている2台のカメラ（視点用のカメラ，ゴミ吸い込み用のカメラ）の映像を配信・MLでのゴミ吸い込み判定を行うサーバーです．


## Required

- Go 1.12 or higher
- [OpenCV](https://opencv.org/)
- [LIBSVM](https://www.csie.ntu.edu.tw/~cjlin/libsvm/)



## Getting Started

```shell
go run server.go <config file>
```

 `<config file>` はデフォルトでは `config/macbook-pro.yaml` が使われます．


## Hardware
この作品のハードウェアについてです。主に以下のハードウェア・部品が必要となります。
- ラジコンRCカー
- 視点用カメラ(Webカメラ)
- ゴミ検知カメラ(RaspberryPi カメラモジュール)
- ペットボトル掃除機

ゴミ検知カメラを使用するためには、市販の掃除機では対応できないため、透明なペットボトルを用いて自作する必要があります。
自作掃除機については以下のリンクを参考に作成してみてください。
[ペットボトルサイクロン掃除機について](http://mikamimikami.hatenablog.com/entry/2019/12/23/001000)

###  Creators:

- [Taiga Mikami](https://taigamikami.netlify.com/): `hardware`, `presentation`
- [Yuki Nakahira](https://raahii.github.io/about/): `software`

