"use strict"
import React, { useEffect, useState } from "react"
import ReactDOM from 'react-dom'
import PropTypes from 'prop-types';

import MessagePostForm from './messageForm';

type User = {
  id: number,
  name: string,
};

type Message = {
  addedBy: User,
  id: number,
  message: string,
};

type GetMessageResult = Message[];

function GetMessage(props) {
  //このコードではページが表示された時点でこの関数が走るのでuseStateの初期値を与えないと undefined error
  //https://www.debuggr.io/react-map-of-undefined/
  //非同期の結果を残すためのもの？
  //reactは非同期得意ではなかったらしいが, react hook　が出てきて改善したらしい
  const [data, setData] = useState<GetMessageResult>([]);

  useEffect(() => {
    // ルート /message に対して GETリクエストを送る
    // 帰ってきたものをjsonにしてuseStateに突っ込む
    fetch("/message")
      .then(res => res.json())
      .then(setData)
      .catch(console.log);
  }, [props.forceUpdate]);


  return (
    //タグが複数できる場合は何らかのタグで全体を囲う
    <div>
      {data.map((item) => (
        //{}で囲むと変数展開できる
        //djangoのtemplateとかもそうだった　流行ってるんかな 便利やし
        <p key={item.id}>{item.id}:[{item.addedBy.name}:{item.addedBy.id}]:{item.message}</p>
      ))}
    </div>
  );
}

GetMessage.propTypes = {
  // このランダム値を変更することで、強制的にサーバーからメッセージを取得させ、最新の情報を入手させる
  forceUpdate: PropTypes.number,
};

GetMessage.defaultProps = {
  forceUpdate: Math.random(),
};

function MessageSection() {
  const [randomValue, setRandomValue] = useState<number>(Math.random());

  const onMessageAdded = () => {
    // フォームによってメッセージが追加されたら、メッセージ一覧を更新する
    setRandomValue(Math.random());
  };

  return (
    <>
      <GetMessage forceUpdate={randomValue} />
      <MessagePostForm onSubmitSuccessful={onMessageAdded} />
    </>
  )
}


//webpackでバンドルしている関係で存在していないIDが指定される場合がある
//エラーをそのままにしておくと、エラー以後のレンダリングがされない
if(document.getElementById('message') != null){
  ReactDOM.render(<MessageSection />, document.getElementById('message'));
}
