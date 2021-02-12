# Overview
Syo is a small program that downloads web novels from Syosetu

```
xxx@xxx /home> ./syo --help
Usage of \xxx\xxx\syo:
  -l duration
        limit requests to syosetu (default 200ms)
  -n string
        ids of novels
  -o    overwrite existing content
  -w int
        number of workers (default 1)
  -s string
        which syosetu site (default is ncode)
```

```
# -l Rate Limiting: 1 second between requests
# -w Workers: 10
# -n Download Novels: n7104eb,n0628ew,n2933fj
# -o Overwrite existing chapters, toc, etc...
# -s Syosetu site: ncode/yomou/novel18/etc.

--will search on https://ncode.syosetu.com
./syo -l 1s -w 10 -o -n n7104eb,n0628ew,n2933fj

--will search on https://novel18.syosetu.com
./syo -l 1s -w 10 -o -s novel18 -n n8641dj
```

Running the above will end up looking something like this for each web novel. the source file contains a link back to the novel url. toc contains a table of contents of the novel chapters and the rest should be self explanatory.

```
wn
    ├── n0628ew
    │   ├── 1.txt
    │   ├── 10.txt
    │   ├── 11.txt
    │   ├── 12.txt
    │   ├── 13.txt
    │   ├── 14.txt
    │   ├── 15.txt
    │   ├── 16.txt
    │   ├── 17.txt
    │   ├── 18.txt
    │   ├── 19.txt
    │   ├── 2.txt
    │   ├── 20.txt
    │   ├── 21.txt
    │   ├── 22.txt
    │   ├── 23.txt
    │   ├── 24.txt
    │   ├── 25.txt
    │   ├── 26.txt
    │   ├── 27.txt
    │   ├── 28.txt
    │   ├── 29.txt
    .   .
    .   .
    .   .
    │   ├── desc.txt
    │   ├── source.txt
    │   ├── title.txt
    │   └── toc.csv
    .
    .
    .
```

toc.csv:
```csv
Timestamp, Revised, URL, Title
2018/07/04 22:45, yes, https://ncode.syosetu.com/n0628ew/1/, 1-1 プロローグ
2018/07/04 23:39, yes, https://ncode.syosetu.com/n0628ew/2/, 1-2　蓋をしようと思います
2018/07/05 22:12, yes, https://ncode.syosetu.com/n0628ew/3/, 1-3　歩き始めようと思います
2018/07/06 23:21, yes, https://ncode.syosetu.com/n0628ew/4/, 2-1　学んでみようと思います
2018/07/07 09:31, yes, https://ncode.syosetu.com/n0628ew/5/, 2-2　目を通してみようと思います
2018/07/07 23:27, yes, https://ncode.syosetu.com/n0628ew/6/, 2-3　どうすべきか考えようと思います
2018/07/08 01:46, yes, https://ncode.syosetu.com/n0628ew/7/, フィルツ＝オルストイ
2018/07/08 23:17, yes, https://ncode.syosetu.com/n0628ew/8/, 2-4　魔導学を学ぼうと思います
2018/07/09 22:17, yes, https://ncode.syosetu.com/n0628ew/9/, 2-5 得意属性を調べてみようと思います
2018/07/11 03:10, yes, https://ncode.syosetu.com/n0628ew/10/, 2-6　全力で励もうと思います
2018/07/12 00:50, yes, https://ncode.syosetu.com/n0628ew/11/, フィルツ＝オルストイ　2
2018/07/12 22:52, yes, https://ncode.syosetu.com/n0628ew/12/, 3-1　王都に行ってみようと思います
2018/07/14 02:57, yes, https://ncode.syosetu.com/n0628ew/13/, 3-2　ご挨拶をしてみようと思います
2018/07/14 11:33, yes, https://ncode.syosetu.com/n0628ew/14/, 3-3　お友達を作ってみようと思います
2018/07/15 08:36, yes, https://ncode.syosetu.com/n0628ew/15/, Pural Bless ～7封剣と光の巫女姫　　アイン＝ファーランド　ルート

...
```

# Build & Run

```
go get -u github.com/SirMetathyst/Syo
cd to/your/go/path/syo
make
```

with [UPX](https://upx.github.io/) 

```
make build_small
```
