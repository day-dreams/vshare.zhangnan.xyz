# vshare.zhangnan.xyz
视频分享/同步观看



- 如何开启新房间
    - 如果要看新视频，先传到腾讯云
    - runtime/config.json，添加相关的room/vid信息
    - 重启之`make daemon`
    - 电脑浏览器访问，url带上query param，`rid=xxxx`即可
- 编译运行，`make docker`

- Reference
    - HLS标准，https://www.rfc-editor.org/pdfrfc/rfc8216.txt.pdf