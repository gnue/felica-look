# felica-look

別ソフトウェアで出力した FeliCaカードのダンプ出力を解析して意味のわかる形式で表示します

## 対応カード

* Suica互換の交通系カード（利用履歴のみ、駅データは入っていません）
* Edy
* RapiCa（鹿児島県の共通乗車IC）

## 対応フォーマット

* felica_dump ([libpafe](http://homepage3.nifty.com/slokar/pasori/libpafe.html)) -- RC-S310, RC-S320, RC-S330
* lpdump ([libpasori](http://sourceforge.jp/projects/libpasori/)) -- RC-S320のみ
* [FeliCa Raw Viewer](http://oasis.halfmoon.jp/mw/index.php?title=Soft-FelicaRawViewer) -- Windowsアプリ

## インストール

1. [Releases](https://github.com/gnue/felica-look/releases) からダウンロード
2. zipファイルを解凍
3. 実行ファイルをターミナルで実行
   * 直接ファイルを指定して実行
   * 実行パスを設定
   * 実行パスのあるディレクトリにシンボリックリンクを作成（おすすめ）

## 使い方

	$ felica-look ファイル        # フォーマット解析して情報を表示
	$ felica-look -e ファイル     # [利用履歴（元データ）] も表示
	$ felica-look -x ファイル     # 元データの16進表示もいっしょに表示します
	$ felica-look -d ファイル     # 内部形式に変換したデータをダンプ表示（デバッグ用）

パイプライン（FeliCaカードを読んですぐに表示する）

	$ felica_dump | felica-look

## YAML

表示用の変換テーブルを YAML形式で持ちます

* 駅名・停留所の名前はここで追加できます

### YAMLファイルの検索ディレクトリ

1. カレントディレクトリ
2. 実行ファイルあるディレクトリ（シンボリックリンクを辿って実行ファイルのあるディレクトリを探します）
3. 実行ファイルあるディレクトリのカードに対応するサブディレクトリ（例：`suica.yml`の場合は`suica/suica.yml`）

## ビルド

ソースコードからビルドする場合

	$ git clone https://github.com/gnue/felica-look
	$ cd felica-look
	$ go build

注：あらかじめ Go言語をインストールしておく必要があります

	$ brew install go

## FAQ

0. なぜ Go言語で開発したのですか？
   * 環境に依存せず配布バイナリのみで使えるようできる
   * また、マルチプラットフォームのバイナリも容易に作成できる
   * C言語のコードを容易に組込める
   * 単に Go言語を使ってみたかったから（初プログラム）
0. なぜ C言語がまざっているのですか？
   * 他ソフトウェアと共通化をはかるため
   * 実際に `suica.h`, `edy.h` は PasoriKit からの流用
   * `rapica.h` は PasoriKit に組込み予定
   * たぶん Go言語だけで書いたほうがシンプルになるのかもしれないけど、同じ構造体定義やロジックをC言語とGo言語で毎回書き直しはしたくないので
0. Suicaの駅名が `suica.yml` に登録されていないのはどうでしてですか？
   * 駅の数が多いため
   * `suica.yml` の `STATIONS:` に登録すれば表示することができます（ただし、表示がずれてしまうので注意）
   * [IC SFCard Fan](http://www014.upp.so-net.ne.jp/SFCardFan/)さんの[サイバネ駅コード調査データベース](http://www.denno.net/SFCardFan/index.php)との対応は `((地区コード<<16) + (線区コード<<8) + 駅順コード)` になります
0. 駅名を `suica.yml` に登録しましたが表示がずれてしまいました
   * 16進表示の英数字も日本語もどちらも１文字として `fmt.Sprintf` でカウントされてしまうために日本語が入ると表示がずれてしまいます
   * うまい解決方法を募集中...
0. カードリーダーからFeliCaカードのデータを直接読込まないのですか？
   * パイプラインを使えばそれに近いことが可能です
   * 他のライブラリや環境依存にはしたくないのでカードリーダーからの直接読込みを行う予定はありません
0. 表示結果をもっとわかりやすくしたい
   * もともとが FeliCaカードのフォーマット確認を目的にしています
   * 目的がずれてしまうので本プログラム内で必要以上の表示整形は行いません
   * さらなる表示のためには機能追加を予定している LTSV/JSON出力を利用し別プログラムで行うのがいいでしょう
     * 例：JSON出力を使ってWebアプリでブラウザ表示
     * 例：LTSV出力を使って [fluentd](http://fluentd.org) で利用履歴の蓄積

## TODO

* JSON出力
* [LTSV](http://ltsv.org)出力（利用履歴のみ）
* その他カードの対応
