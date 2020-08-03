import React, { useState } from 'react';
import PropTypes from 'prop-types';
import ReactDOM from 'react-dom';
import { createStyles, makeStyles, Theme } from '@material-ui/core/styles';
import CircularProgress from '@material-ui/core/CircularProgress';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import Alert from '@material-ui/lab/Alert';

// メッセージ追加のAPIへのURL
// eslint-disable-next-line @typescript-eslint/naming-convention
const API_URL_LOGIN = '/register';

function RegisterPostForm(props) {
const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      '& > *': {
        margin: theme.spacing(1),
        width: '25ch',
      },
      width: '100%',
      '& > * + *': {
        marginTop: theme.spacing(2),
      },
    },
  }),
);

  // テキストボックス内のメッセージ
  const [userId, setUserId] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  // サーバがへメッセージ追加のリクエストを処理中ならtrue、でないならfalseの状態
  const [working, setWorking] = useState<boolean>(false);

  const classes = useStyles();
  const [progress, setProgress] = React.useState(0);

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
        ReactDOM.render(<div className={classes.root}><Alert variant="outlined" severity="error" onClose={() => {
          ReactDOM.render(<div></div>, document.getElementById('serverMessage'));
        }}>{obj.error}</Alert></div>, document.getElementById('serverMessage'));
        throw new Error(`An error occurred on querying ${API_URL_LOGIN}, the response included error message: ${obj.error}`);
      }
      if (!('success' in obj)) {
        // サーバーからsuccessメンバが含まれたJSONが帰るはずだが、見当たらなかった
        ReactDOM.render(<div className={classes.root}><Alert variant="outlined" severity="error" onClose={() => {
          ReactDOM.render(<div></div>, document.getElementById('serverMessage'));
        }}>Error</Alert></div>, document.getElementById('serverMessage'));
        throw new Error(`An response from ${API_URL_LOGIN} unexpectedly did not have 'success' member`);
      }
      if (obj.success !== true) {
        ReactDOM.render(<div className={classes.root}><Alert variant="outlined" severity="error" onClose={() => {
          ReactDOM.render(<div></div>, document.getElementById('serverMessage'));
        }}>Error</Alert></div>, document.getElementById('serverMessage'));
        throw new Error(`An response from ${API_URL_LOGIN} returned non true value as 'success' member`);
      }

      // 要求は成功
      ReactDOM.render(<div className={classes.root}>
        <Alert variant="outlined" severity="success">登録完了! 3秒後にログインページへ推移</Alert>
        <CircularProgress />
        </div>, document.getElementById('serverMessage'));
      // リスナ関数を呼ぶ
      props.onSubmitSuccessful();

      // 登録が成功したらログインページにリダイレクト
      setTimeout(() => {
        window.location.href = '/login';
      }, 3000);
    } finally {
      setWorking(false);
    }
  };

  return (
    <>
      <form className={classes.root} noValidate autoComplete="off">
      <TextField id="standard-basic" label="ユーザーID" value={userId}
      type='textbox'
      onChange={(event) => setUserId(event.target.value)}/>
      <p></p>
      <TextField id="standard-basic" label="パスワード" value={password}
        type='textbox'
        onChange={(event) => setPassword(event.target.value)}/>
      <p></p>
      <Button variant="contained" color="primary"　disabled={working} onClick={handleSubmit}>
        登録
      </Button>
      </form>
    </>
  );
}

RegisterPostForm.propTypes = {
  onSubmitSuccessful: PropTypes.func,
};

RegisterPostForm.defaultProps = {
  onSubmitSuccessful: () => {},
};

//
// //webpackでバンドルしている関係で存在していないIDが指定される場合がある
// //エラーをそのままにしておくと、エラー以後のレンダリングがされない
// if(document.getElementById('register') != null){
//   ReactDOM.render(<RegisterPostForm />, document.getElementById('register'));
// }
