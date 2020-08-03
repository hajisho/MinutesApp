import React, { useState } from 'react';
import PropTypes from 'prop-types';

import ReactDOM from 'react-dom';
// メッセージ追加のAPIへのURL
// eslint-disable-next-line @typescript-eslint/naming-convention
const API_URL_LOGIN = '/login';

function LoginPostForm(props) {
  // テキストボックス内のメッセージ
  const [userId, setUserId] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  // サーバがへメッセージ追加のリクエストを処理中ならtrue、でないならfalseの状態
  const [working, setWorking] = useState<boolean>(false);

  const handleSubmit = async (event: React.FormEvent) => {
    // FIXME もしかしたら、非同期なため、これが効く前にボタンをクリックできるかもしれない
    setWorking(true);
    try {
      // ページが更新されないようにする
      event.preventDefault();

      // Reactのハンドラはasyncにできる
      const res = await fetch(API_URL_LOGIN, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        // 相応しくないかも
        // same-originを使うべき？
        credentials: 'include',
        body: JSON.stringify({ userId, password }),
      });

      setUserId('');
      setPassword('');

      const obj = await res.json();
      if ('error' in obj) {
        // サーバーからエラーが返却された
        ReactDOM.render(
          <p>{obj.error}</p>,
          document.getElementById('serverMessage')
        );
        throw new Error(
          `An error occurred on querying ${API_URL_LOGIN}, the response included error message: ${obj.error}`
        );
      }
      if (!('success' in obj)) {
        // サーバーからsuccessメンバが含まれたJSONが帰るはずだが、見当たらなかった
        ReactDOM.render(<p>error</p>, document.getElementById('serverMessage'));
        throw new Error(
          `An response from ${API_URL_LOGIN} unexpectedly did not have 'success' member`
        );
      }
      if (obj.success !== true) {
        ReactDOM.render(<p>error</p>, document.getElementById('serverMessage'));
        throw new Error(
          `An response from ${API_URL_LOGIN} returned non true value as 'success' member`
        );
      }
      // 要求は成功
      ReactDOM.render(
        <p>ログイン成功!</p>,
        document.getElementById('serverMessage')
      );
      // リスナ関数を呼ぶ
      props.onSubmitSuccessful();

      // ログインが成功したらmainページにリダイレクト
      window.location.href = '/';
    } finally {
      setWorking(false);
    }
  };

  return (
    <>
      <a href="/register">
        <button type="button">登録ページへ</button>
      </a>
      <form onSubmit={handleSubmit}>
        <input
          value={userId}
          type="textbox"
          placeholder="ユーザーID"
          onChange={(event) => setUserId(event.target.value)}
        />
        <input
          value={password}
          type="textbox"
          placeholder="パスワード"
          onChange={(event) => setPassword(event.target.value)}
        />
        <button type="submit" disabled={working}>
          ログイン
        </button>
      </form>
    </>
  );
}

LoginPostForm.propTypes = {
  onSubmitSuccessful: PropTypes.func,
};

LoginPostForm.defaultProps = {
  onSubmitSuccessful: () => {},
};

// webpackでバンドルしている関係で存在していないIDが指定される場合がある
// エラーをそのままにしておくと、エラー以後のレンダリングがされない
if (document.getElementById('login') != null) {
  ReactDOM.render(<LoginPostForm />, document.getElementById('login'));
}
