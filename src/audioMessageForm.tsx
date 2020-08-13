import React, { useState } from 'react';
import Icon from '@material-ui/core/Icon';
import Button from '@material-ui/core/Button';
import PropTypes from 'prop-types';

const apiUrlAddMessage = '/add_message';

export default function AudioMessagePostForm(props) {
  const [working, setWorking] = useState<boolean>(false);

  if ('SpeechRecognition' in window) {
    (window as any).SpeechRecognition = (window as any).SpeechRecognition;
  } else if ('webkitSpeechRecognition' in window) {
    (window as any).SpeechRecognition = (window as any).webkitSpeechRecognition;
  }

  const speech = new (window as any).SpeechRecognition();
  speech.lang = 'ja-JP';
  speech.onresult = async function AudioResult(e) {
    speech.stop();
    if (e.results[0].isFinal) {
      const audioText = e.results[0][0].transcript;
      const res = await fetch(apiUrlAddMessage, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ message: audioText }),
      });
      const obj = await res.json();
      if ('error' in obj) {
        // サーバーからエラーが返却された
        throw new Error(
          `An error occurred on querying ${apiUrlAddMessage}, the response included error message: ${obj.error}`
        );
      }
      if (!('success' in obj)) {
        // サーバーからsuccessメンバが含まれたJSONが帰るはずだが、見当たらなかった
        throw new Error(
          `An response from ${apiUrlAddMessage} unexpectedly did not have 'success' member`
        );
      }
      if (obj.success !== true) {
        throw new Error(
          `An response from ${apiUrlAddMessage} returned non true value as 'success' member`
        );
      }
      props.onSubmitSuccessful();
    }
  };

  speech.onend = () => {
    speech.start();
  };

  const handleSubmit = (event: React.FormEvent) => {
    setWorking(true);
    try {
      // ページが更新されないようにする
      speech.start();
      event.preventDefault();
      props.onSubmitSuccessful();
    } finally {
      setWorking(false);
    }
  };
  return (
    <form>
      <Button
        disabled={working}
        variant="contained"
        color="primary"
        endIcon={<Icon>send</Icon>}
        onClick={handleSubmit}
      >
        Start
      </Button>
      <Button
        disabled={working}
        variant="contained"
        color="primary"
        onClick={() => {
          window.location.href = '/';
        }}
      >
        Stop
      </Button>
    </form>
  );
}
AudioMessagePostForm.propTypes = {
  // 新しいメッセージの追加が正常に完了したら呼ばれる関数
  onSubmitSuccessful: PropTypes.func,
};

AudioMessagePostForm.defaultProps = {
  onSubmitSuccessful: () => {},
};
