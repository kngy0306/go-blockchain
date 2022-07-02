# go-blockchain

## 現役シリコンバレーエンジニアが教えるGoで始めるスクラッチからのブロックチェーン開発入門

https://www.udemy.com/course/go-blockchain/


### ブロックチェーン

ブロックチェーンで利用するサーバは誰が用意するか？
  個人がサーバを提供している。参加者はMiningを行い報酬を受け取ることができる。

### セクション3

Blockに含むもの

- prev hash
  - 前のBlockの情報を用いて作成したハッシュ
- timestamp
  - 生成時間
- nonce
  - ルールに基づいた解答？
- transactions
  - 取引時の情報
  - 実際はPoolと呼ばれる領域にtransactionを記録し、マイニングが終わればtransactionsに書き込む

BlockChainに含むもの

- transaction pool
  - Transaction構造体をポインタでもつ
- chain
  - Block構造体をポインタでもつ

nonceについて
新しいブロックが正しいものと判断するため、コンセンサスアルゴリズムを用いて解を求める。解を計算して、nonceが正しいと判断する。

- challenge値
- prev hash
- transaction

を用いる。challenge値がnonceの候補値になる。ハッシュ値の先頭が0の3桁となれば、そのchallenge値がnonceとなる。0の桁数は難易度設定により変動する。本講座は0を3桁を解答とする。