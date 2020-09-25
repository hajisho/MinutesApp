import React, { useState } from 'react';
import PropTypes from 'prop-types';
import GetAppIcon from '@material-ui/icons/GetApp';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from '@material-ui/core';

export default function DownloadMessageDialog(props) {
  const { targetMessage, title } = props;

  // Dialogを開くかどうか
  const [open, setOpen] = useState<boolean>(false);

  const handleDownload = () => {
    const blob = new Blob([targetMessage], { type: 'text/plan' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `${title}.txt`;
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
      <Button
        onClick={handleClickOpen}
        size="small"
        color="inherit"
        endIcon={<GetAppIcon fontSize="small" />}
      >
        Download
      </Button>
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
  title: PropTypes.string,
};

DownloadMessageDialog.defaultProps = {
  targetMessage: '',
  title: '',
};
