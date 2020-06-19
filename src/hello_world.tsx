"use strict"
import React, { useEffect, useState } from "react"
import ReactDOM from 'react-dom'

function GetMessage() {
  //このコードではページが表示された時点でこの関数が走るのでuseStateの初期値を与えないと undefined error
  //https://www.debuggr.io/react-map-of-undefined/
  //非同期の結果を残すためのもの？
  //reactは非同期得意ではなかったらしいが, react hook　が出てきて改善したらしい
  const [data, setData] = useState([]);

  useEffect(() => {
    //直でasync関数を受け取れないので一度噛ませる
    const res = async() => {
      //ルート/message　に対して GETリクエストを送る
      //帰ってきたものをjsonにしてuseStateに突っ込む
      const r = await fetch("/message")
                          .then(res => res.json())
                          .then(setData)
    }
    res();

  }, []);

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

ReactDOM.render(<GetMessage />, document.getElementById('message'));

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
