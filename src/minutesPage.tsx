import React, { useState, useEffect } from 'react';
import ReactDOM from 'react-dom';
import PropTypes from 'prop-types';
import { makeStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import Menu from '@material-ui/core/Menu';
import AccountCircle from '@material-ui/icons/AccountCircle';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import Avatar from '@material-ui/core/Avatar';

import MessagePostForm from './messageForm';
import EditMessagePostForm from './editForm';

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
function GetMessage(props) {
  const { forceUpdate } = props;
  const classes = useStylesCard();

  type User = {
    id: number;
    name: string;
  };

  type Message = {
    addedBy: User;
    id: number;
    message: string;
  };

  type GetMessageResult = Message[];
  const [data, setData] = useState<GetMessageResult>([]);

  useEffect(() => {
    // ルート /message に対して GETリクエストを送る
    // 帰ってきたものをjsonにしてuseStateに突っ込む
    fetch('/message')
      .then((res) => res.json())
      .then(setData);
  }, [forceUpdate]);

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

function MessageSection() {
  const [randomValue, setRandomValue] = useState<number>(Math.random());

  const onMessageAdded = () => {
    // フォームによってメッセージが追加されたら、メッセージ一覧を更新する
    setRandomValue(Math.random());
  };

  return (
    <>
      <MessagePostForm onSubmitSuccessful={onMessageAdded} />
      <GetMessage forceUpdate={randomValue} />
    </>
  );
}

const useStylesBar = makeStyles((theme) => ({
  header: {
    flexGrow: 1,
  },
  menuButton: {
    marginRight: theme.spacing(2),
  },
  title: {
    flexGrow: 1,
  },
}));

function MinuteAppBar() {
  const classes = useStylesBar();
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

  const isMenuOpen = Boolean(anchorEl);

  const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const menuId = 'primary-search-account-menu';
  const renderMenu = (
    <Menu
      anchorEl={anchorEl}
      anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
      id={menuId}
      keepMounted
      transformOrigin={{ vertical: 'top', horizontal: 'right' }}
      open={isMenuOpen}
      onClose={handleMenuClose}
    >
      <MenuItem onClick={handleMenuClose}>Profile</MenuItem>
      <MenuItem onClick={handleMenuClose}>My account</MenuItem>
    </Menu>
  );

  return (
    <div className={classes.header}>
      <AppBar position="static">
        <Toolbar>
          <IconButton
            edge="start"
            className={classes.menuButton}
            color="inherit"
            aria-label="menu"
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" className={classes.title}>
            Minutes Application
          </Typography>
          <IconButton
            edge="end"
            aria-label="account of current user"
            aria-controls={menuId}
            aria-haspopup="true"
            onClick={handleProfileMenuOpen}
            color="inherit"
          >
            <AccountCircle />
          </IconButton>
        </Toolbar>
      </AppBar>
      {renderMenu}
    </div>
  );
}

// webpackでバンドルしている関係で存在していないIDが指定される場合がある
// エラーをそのままにしておくと、エラー以後のレンダリングがされない
if (document.getElementById('message') != null) {
  ReactDOM.render(<MessageSection />, document.getElementById('message'));
}
if (document.getElementById('minuteHeader') != null) {
  ReactDOM.render(<MinuteAppBar />, document.getElementById('minuteHeader'));
}
