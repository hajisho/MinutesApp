"use strict"
import React, { useEffect, useState } from "react"
import ReactDOM from 'react-dom'
import PropTypes from 'prop-types';

import MessagePostForm from './form';

type Message = {
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

  console.log(data)

  return (
    //タグが複数できる場合は何らかのタグで全体を囲う
    <div>
      {data.map((item) => (
        //{}で囲むと変数展開できる
        //djangoのtemplateとかもそうだった　流行ってるんかな 便利やし
        <p>{item.id}:{item.message}</p>
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

ReactDOM.render(<MessageSection />, document.getElementById('message'));

//文字通り
//理解が進んだら消す予定
class Hello extends React.Component{
  render(){
    return(
      <div>Hello World</div>
    );
  }
}

ReactDOM.render(<Hello />, document.getElementById("content"))
