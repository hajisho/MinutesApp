import React, { useState } from 'react';
import PropTypes from 'prop-types';
import GetAppIcon from '@material-ui/icons/GetApp';
import IconButton from '@material-ui/core/IconButton';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';

export default function DownloadMessageDialog(props) {
  const { targetMessage /* title */ } = props;

  // Dialogを開くかどうか
  const [open, setOpen] = useState<boolean>(false);

  const handleDownload = () => {
    const blob = new Blob([targetMessage], { type: 'text/plan' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = 'aaa.txt'; // `${title}.txt`;
    link.click();

    setOpen(false);
  };

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <span>
      <IconButton
        aria-label="download"
        onClick={handleClickOpen}
        data-num="100"
      >
        <GetAppIcon fontSize="small" />
      </IconButton>
      <Dialog
        open={open}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">
          この会議をダウンロードしますか？
        </DialogTitle>
        <DialogContent>
          <DialogContentText id="download-dialog-description">
            {targetMessage}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} color="inherit" size="small" autoFocus>
            Cancel
          </Button>
          <Button
            onClick={handleDownload}
            color="primary"
            size="small"
            endIcon={<GetAppIcon fontSize="small" />}
            variant="contained"
          >
            Download
          </Button>
        </DialogActions>
      </Dialog>
    </span>
  );
}

DownloadMessageDialog.propTypes = {
  targetMessage: PropTypes.string,
  // title: PropTypes.string,
};

DownloadMessageDialog.defaultProps = {
  targetMessage: '',
  // title: '',
};
