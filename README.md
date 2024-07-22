# koron/c3tr-client

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron/c3tr-client)](https://pkg.go.dev/github.com/koron/c3tr-client)
[![Actions/Go](https://github.com/koron/c3tr-client/workflows/Go/badge.svg)](https://github.com/koron/c3tr-client/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron/c3tr-client)](https://goreportcard.com/report/github.com/koron/c3tr-client)

[llama.cpp](https://github.com/ggerganov/llama.cpp) の llama-server でローカルで動かしてる [C3TR-Adapter\_gguf](https://huggingface.co/webbigdata/C3TR-Adapter_gguf) にAPI (`/completions`) 経由で翻訳させるクライアントプログラム。

## Gettings Started

1. Install [llama.cpp](https://github.com/ggerganov/llama.cpp/releases/latest)

    Windows + CUDA の場合は `llama-b{数字}-win-cuda-cu{CUDAのバージョン}-x64.zip` の何れかをダウンロードして展開する。

2. Download [C3TR-Adapter\_gguf](https://huggingface.co/webbigdata/C3TR-Adapter_gguf/tree/main)

    `C3TR-Adapter-Q4_k_m.gguf` もしくは `C3TR-Adapter.f16.Q4_k_m.gguf` あたりがオススメ。

3. c3tr-client をインストールする

    ```console
    $ go install github.com/koron/c3tr-client@latest
    ```

4. (OPTIONAL) Setup environment variables

    CUDA用とllama-server用の環境変数を設定する。

    以下は筆者の設定例:

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

5. c3tr-client を使って翻訳する

    ```console
    $ c3tr-client "A client for the C3TR Agent for Japanese-English and English-Japanese translation running on llama.cpp"
    llama.cpp上で動作するC3TRエージェントの日本語-英語と英語-日本語の翻訳クライアント

    $ c3tr-client "llama.cpp上で動作するC3TRエージェントの日本語-英語と英語-日本語の翻訳クライアント"
    A Japanese-English and English-Japanese translation client for the C3TR agent that runs on llama.cpp
    ```
