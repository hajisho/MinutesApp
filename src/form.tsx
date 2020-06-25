import React, { useState } from 'react';
import PropTypes from 'prop-types';

// メッセージ追加のAPIへのURL
const API_URL_ADD_MESSAGE = '/add_message';

export default function MessagePostForm(props) {
  // テキストボックス内のメッセージ
  const [message, setMessage] = useState<string>('');
  // サーバがへメッセージ追加のリクエストを処理中ならtrue、でないならfalseの状態
  const [working, setWorking] = useState<boolean>(false);

  const handleSubmit = async (event: React.FormEvent) => {
    // FIXME もしかしたら、非同期なため、これが効く前にボタンをクリックできるかもしれない
    setWorking(true)
    try {
      // ページが更新されないようにする
      event.preventDefault();

      // Reactのハンドラはasyncにできる
      const res = await fetch(API_URL_ADD_MESSAGE, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ message }),
      });
      const obj = await res.json();
      if ('error' in obj) {
        // サーバーからエラーが返却された
        throw new Error(`An error occurred on querying ${API_URL_ADD_MESSAGE}, the response included error message: ${obj.error}`);
      }
      if (!('success' in obj)) {
        // サーバーからsuccessメンバが含まれたJSONが帰るはずだが、見当たらなかった
        throw new Error(`An response from ${API_URL_ADD_MESSAGE} unexpectedly did not have 'success' member`);
      }
      if (obj.success !== true) {
        throw new Error(`An response from ${API_URL_ADD_MESSAGE} returned non true value as 'success' member`);
      }
      // 要求は成功
      // リスナ関数を呼ぶ
      props.onSubmitSuccessful();
    } finally {
      setWorking(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type='textbox'
        placeholder='ここに追加したいメッセージを入力します'
        onChange={(event) => setMessage(event.target.value)}
      />
      <button disabled={working}>追加</button>
    </form>
  )
}

MessagePostForm.propTypes = {
  // 新しいメッセージの追加が正常に完了したら呼ばれる関数
  onSubmitSuccessful: PropTypes.func,
};

MessagePostForm.defaultProps = {
  onSubmitSuccessful: () => {},
};
