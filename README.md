htn2hugo
===

はてなブログの記事を Hugo 用の Markdown に変換するツールです。

## Features

- はてなブログのすべての記法 (Markdown/はてな記法/見たまま) に対応しています
- はてな記法で書かれた記事の下記の要素は, それぞれ HTML に変換されます
  - リンク文字列
  - はてなブログカード
  - Twitter 埋め込み
  - はてなフォトライフの画像埋め込み
- 下書きの記事は `draft: true` として生成されます
- Permalink は `url: "/entry/..."` として生成されます 

## Example

はてなブログの AtomPub で取得できた XML の中に下記のような `entry` が含まれていた場合の変換例です。

```xml
<entry>
<id>tag:blog.hatena.ne.jp,2013:blog-michimani-12921228815716705189-26006619999999999</id>
<link rel="edit" href="https://blog.hatena.ne.jp/michimani/michimani.hateblo.jp/atom/entry/26006619999999999"/>
<link rel="alternate" type="text/html" href="https://www.michinoeki-mania.com/entry/2020/12/22/999999"/>
<author><name>michimani</name></author>
<title>test</title>
<updated>2020-12-22T00:03:46+09:00</updated>
<published>2020-12-21T21:00:58+09:00</published>
<app:edited>2020-12-22T00:03:46+09:00</app:edited>
<summary type="text">header 1 Thie is a sample post. header 2 はてな記法という罠— よっしー Lv.859 | michimani (@michimani210) 2020年12月20日 $ echo &quot;Hello World!&quot;</summary>
<content type="text/x-markdown"># header 1

Thie is a sample post.

## header 2

[f:id:michimani:20180912345678j:plain]

[https://twitter.com/michimani210/status/1340486979226460160:embed]


```bash
$ echo &quot;Hello World!&quot;
```</content>
<hatena:formatted-content type="text/html" xmlns:hatena="http://www.hatena.ne.jp/info/xmlns#">&lt;h1&gt;header 1&lt;/h1&gt;

&lt;p&gt;Thie is a sample post.&lt;/p&gt;

&lt;h2&gt;header 2&lt;/h2&gt;

&lt;p&gt;&lt;span itemscope itemtype=&quot;http://schema.org/Photograph&quot;&gt;&lt;img src=&quot;https://cdn-ak.f.st-hatena.com/images/fotolife/m/michimani/20180911/20180912345678.jpg&quot; alt=&quot;f:id:michimani:20180912345678j:plain&quot; title=&quot;&quot; class=&quot;hatena-fotolife&quot; itemprop=&quot;image&quot;&gt;&lt;/span&gt;&lt;/p&gt;

&lt;p&gt;&lt;blockquote data-conversation=&quot;none&quot; class=&quot;twitter-tweet&quot; data-lang=&quot;ja&quot;&gt;&lt;p lang=&quot;ja&quot; dir=&quot;ltr&quot;&gt;はてな記法という罠&lt;/p&gt;&amp;mdash; よっしー Lv.859 | michimani (@michimani210) &lt;a href=&quot;https://twitter.com/michimani210/status/1340486979226460160?ref_src=twsrc%5Etfw&quot;&gt;2020年12月20日&lt;/a&gt;&lt;/blockquote&gt; &lt;script async src=&quot;https://platform.twitter.com/widgets.js&quot; charset=&quot;utf-8&quot;&gt;&lt;/script&gt; &lt;/p&gt;

&lt;pre class=&quot;code bash&quot; data-lang=&quot;bash&quot; data-unlink&gt;$ echo &amp;#34;Hello World!&amp;#34;&lt;/pre&gt;

</hatena:formatted-content>

<category term="作ってみました" />

<category term="その他" />

<app:control>
  <app:draft>yes</app:draft>
</app:control>

</entry>
```

#### 変換後の Markdown

```markdown
---
title: "test"
date: 2020-12-22T00:03:46+09:00
draft: true
author: ["michimani"]
categories: ["作ってみました","その他"]
archives: ["2020", "2020-12"]
description: "header 1 Thie is a sample post. header 2 はてな記法という罠— よっしー Lv.859 | michimani (@michimani210) 2020年12月20日 $ echo \"Hello World!\""
url: "/entry/2020/12/22/999999"
---

# header 1

Thie is a sample post.

## header 2


<a href="https://f.hatena.ne.jp/michimani/20180912345678">
  <img src="https://cdn-ak.f.st-hatena.com/images/fotolife/m/michimani/20180911/20180912345678.jpg" alt="20180912345678">
</a>


<blockquote class="twitter-tweet" >
  <p lang="ja" dir="ltr"></p>
  <a href="https://twitter.com/michimani210/status/1340486979226460160"></a>
</blockquote>
<script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>



```bash
$ echo "Hello World!"
\```

```

※ 最後のコードブロックの閉じに `\` がついていますが、本来は付きません

## Usage

1. clone repository

```bash
$ git clone https://github.com/michimani/htn2hugo.git
```

2. build docker image

```bash
$ cd htn2hugo
$ docker build -t htn2hugo .
```

3. run docker container

実行時に下記の環境変数を渡します。

- `HTN_ID` : はてなの ID です。
- `HTN_API_KEY` : はてなブログの API キーです。はてなブログの管理画面から **設定 > 詳細設定** で確認できます。

```bash
$ docker run \
-e HTN_ID=<your hatena ID> \
-e HTN_API_KEY=<your hatena API key> \
-v <path to your hugo dir>/content/posts:/dist htn2hugo
```

## Note

一部 HTML として変換するため、 Hugo の `config.toml` 内では下記の設定が必要です。

```toml
[markup]
    [markup.goldmark]
        [markup.goldmark.renderer]
            unsafe = true
```