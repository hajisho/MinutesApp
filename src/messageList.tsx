import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import {
  makeStyles,
  Card,
  CardContent,
  Avatar,
  Typography,
  CardHeader,
  CardActions,
} from '@material-ui/core';
// eslint-disable-next-line no-unused-vars
import { Message, User } from './datatypes';
import EditMessagePostForm from './editForm';
import DeleteMessageDialog from './deleteDialog';
import AudioMessagePostForm from './audioMessageForm';

const useStylesCard = makeStyles({
  root: {
    minWidth: 275,
    maxWidth: 275,
    marginTop: 15,
    marginBottom: 15,
  },
  bullet: {
    display: 'inline-block',
    margin: '0 2px',
    transform: 'scale(0.8)',
  },
  title: {
    fontSize: 14,
  },
  pos: {
    marginBottom: 12,
  },
});

function GetMessage({ forceUpdate }) {
  const classes = useStylesCard();
  const [data, setData] = useState<Message[]>([]);
  // ユーザー情報を取得
  const [userData, setUserData] = useState<User>({ id: 0, name: '' });

  useEffect(() => {
    // ルート /message に対して GETリクエストを送る
    // 帰ってきたものをjsonにしてuseStateに突っ込む
    fetch('/message')
      .then((res) => res.json())
      .then(setData);
  }, [forceUpdate]);

  useEffect(() => {
    fetch('/user')
      .then((res) => res.json())
      .then(setUserData);
  }, []);

  return (
    // タグが複数できる場合は何らかのタグで全体を囲う
    <div>
      {data.map((item) => (
        <Card className={classes.root} key={item.id}>
          <CardContent>
            <CardHeader
              avatar={<Avatar>{item.addedBy.name}</Avatar>}
              title={item.addedBy.name}
            />
            <Typography variant="body2" component="p">
              {item.message}
            </Typography>
          </CardContent>
          <CardActions>
            <EditMessagePostForm
              prevMessage={item.message}
              id={item.id.toString()}
              isHidden={userData.id !== item.addedBy.id}
            />
            <DeleteMessageDialog
              targetMessage={item.message}
              id={item.id.toString()}
              isHidden={userData.id !== item.addedBy.id}
            />
          </CardActions>
        </Card>
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

export default function MessageList() {
  const [randomValue, setRandomValue] = useState<number>(Math.random());

  const onMessageAdded = () => {
    // フォームによってメッセージが追加されたら、メッセージ一覧を更新する
    setRandomValue(Math.random());
  };

  return (
    <>
      <AudioMessagePostForm onSubmitSuccessful={onMessageAdded} />
      <GetMessage forceUpdate={randomValue} />
    </>
  );
}
