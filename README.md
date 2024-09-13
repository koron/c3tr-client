# koron/c3tr-client

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron/c3tr-client)](https://pkg.go.dev/github.com/koron/c3tr-client)
[![Actions/Go](https://github.com/koron/c3tr-client/workflows/Go/badge.svg)](https://github.com/koron/c3tr-client/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron/c3tr-client)](https://goreportcard.com/report/github.com/koron/c3tr-client)

[llama.cpp](https://github.com/ggerganov/llama.cpp) の llama-server でローカルで動かしてる [C3TR-Adapter\_gguf](https://huggingface.co/webbigdata/C3TR-Adapter_gguf) にAPI (`/completions`) 経由で翻訳させるクライアントプログラム。

## Gettings Started

1. Install [llama.cpp](https://github.com/ggerganov/llama.cpp/releases/latest)

    Windows + CUDA の場合は `llama-b{数字}-win-cuda-cu{CUDAのバージョン}-x64.zip` の何れかをダウンロードして展開する。

    CPUだけで頑張りたい場合は `llama-b{数字}-bin-win-avx2-x64.zip` をダウンロードして展開する。

    macos は `llama-b{数字}-bin-macos-arm64.zip` をダウンロードして展開する。

    LinuxでUbuntuはコンパイル済みバイナリがあるが、それ以外の場合は自分でコンパイルする必要があるだろう。

2. Download [C3TR-Adapter\_gguf](https://huggingface.co/webbigdata/C3TR-Adapter_gguf/tree/main)

    `C3TR-Adapter-Q4_k_m.gguf` もしくは `C3TR-Adapter.f16.Q4_k_m.gguf` あたりがオススメ。

3. c3tr-client をインストールする

    ```console
    $ go install github.com/koron/c3tr-client@latest
    ```

    もしくはリリースページから[ビルド済みバイナリ](https://github.com/koron/c3tr-client/releases/latest)をダウンロードしてもよい。

4. (OPTIONAL) Setup environment variables

    CUDA用とllama-server用の環境変数を設定する。

    以下は筆者の Windows 11 + CUDA 12 の設定例:

    ```bat
    SET "CUDA_PATH=D:\App\NVIDIA GPU Computing Toolkit\CUDA\v12.2"
    SET "CUDA_PATH_V12_2=%CUDA_PATH%"

    SET "LLAMA_PATH=D:\App\llama\current"

    PATH %CUDA_PATH%\bin;%CUDA_PATH%\libnvvp;%CUDA_PATH%\extras\CUPTI\lib64;%LLAMA_PATH%;%PATH%
    ```

4. llama-server を起動する

    ```
    llama-server --log-disable -m D:\var\llamacpp\C3TR-Adapter-Q4_k_m.gguf -ngl 43
    ```

    ローカル `127.0.0.1:8080` で Open AI の互換APIが動き出す。

5. c3tr-client を使って翻訳する

    ```console
    $ c3tr-client "A client for the C3TR Agent for Japanese-English and English-Japanese translation running on llama.cpp"
    llama.cpp上で動作するC3TRエージェントの日本語-英語と英語-日本語の翻訳クライアント

    $ c3tr-client "llama.cpp上で動作するC3TRエージェントの日本語-英語と英語-日本語の翻訳クライアント"
    A Japanese-English and English-Japanese translation client for the C3TR agent that runs on llama.cpp
    ```

    c3tr-client は 4 で開始した Open AI の互換APIにアクセスして翻訳をする。

## Usage

1. 引数に翻訳する文を指定する

    ```console
    $ c3tr-client '引数に翻訳する文を指定する'
    Specify the sentence to translate in the argument
    ```

2. 引数を指定せずに起動すると、インタラクティブに翻訳する

    ```console
    $ c3tr-client
    c3tr> 引数に翻訳する文を指定する
    Specify the sentence to translate in the argument
    c3tr> 引数を指定せずに起動すると、インタラクティブに翻訳する
    When launched without arguments, it translates interactively.
    c3tr>
    ```

    インタラクティブモードは `<CTRL+D>` で終了できる。

    このモードではいくつかのショートカットで履歴にアクセスできる。
    ショートカットの詳細は [peterh/liner](https://github.com/peterh/liner/blob/v1.2.2/README.md#line-editing) を参照のこと。

    パイプやリダイレクトを付けて起動した場合は、インタラクティブモードでは起動できずエラーになる。

    ```console
    $ c3tr-client > /dev/null
    2024/09/13 12:48:24 no text to translate. for enabling interactive mode, do not use pipes nor redirects
    ```

    `-iteration` オプションはインタラクティブモードでは作用しない。

## Options

* `-verbose` デバッグ用のメッセージを表示する
* `-entrypoint {URL}` 翻訳APIのエントリーポイントを指定する。

    特に指定しなければ `http://127.0.0.1:8080/completions` で、ローカルで動いている llama.cpp を利用する。

* `-mode {MODE}` 日英・英日の翻訳モードを指定する。

    特に指定しない場合は自動判定で、翻訳対象の文中に英数字の文字数が75%を越えたら英→日の翻訳となり、そうでなければ日→英の翻訳となる。
    明示的に指定する場合は `e2j` もしくは `EtoJ` で英→日翻訳、`j2e` もしくは `JtoE` で英→日翻訳になる。
    大文字小文字の区別はしない。

* `-writingstyle {WRITING_STYLE}` 訳出文のスタイルを指定する。

    デフォルトは `technical` 。
    有効な値は次の11通り:
    `casual`, `formal`, `technical`, `journalistic`, `web-fiction`, `business`,
    `nsfw`, `educational-casual`, `academic-presentation`, `slang`,
    `sns-casual`
    ([出展](https://huggingface.co/webbigdata/C3TR-Adapter/discussions/1#669e6ef419d0f96d8a77128b))

* `-iteration {count}` 反復翻訳回数を指定する。反復のたびに翻訳モードは逆転する。

    デフォルトは0で、1回限りの翻訳をする。
    1以上を指定した場合、その回数、翻訳を反復する。
    -1を指定した場合、翻訳文の履歴に一致する文章が訳出されるまで、反復翻訳を繰り返す。
