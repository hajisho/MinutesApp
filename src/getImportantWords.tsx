import React, { useState, useEffect } from 'react';

export default function GetImportantWords() {
  const [data, setData] = useState<String[]>([]);

  useEffect(() => {
    // ルート /message に対して GETリクエストを送る
    // 帰ってきたものをjsonにしてuseStateに突っ込む
    fetch('/important_words')
      .then((res) => res.json())
      .then(setData);
  }, []);

  return (
    // タグが複数できる場合は何らかのタグで全体を囲う
    <div>
      {data.map((item) => (
        <p>{item}</p>
      ))}
    </div>
  );
}
