import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { makeStyles, Card, CardContent, CardHeader } from '@material-ui/core';
import CardActions from '@material-ui/core/CardActions';
import Button from '@material-ui/core/Button';
// eslint-disable-next-line no-unused-vars
import { Meeting } from './datatypes';
import CreateMeetingForm from './createMeetingForm';

const useStylesCard = makeStyles({
  root: {
    minWidth: 275,
    maxWidth: 275,
    marginTop: 15,
    marginBottom: 15,
  },
  title: {
    fontSize: 14,
  },
});

function MeetingList({ forceUpdate }) {
  const classes = useStylesCard();
  const [data, setData] = useState<Meeting[]>([]);

  useEffect(() => {
    fetch('/meetings')
      .then((res) => res.json())
      .then(setData);
  }, [forceUpdate]);

  return (
    <div>
      {data.map((m) => (
        <Card className={classes.root} key={m.name}>
          <CardContent>
            <CardHeader title={m.name} />
            {/* <Typography variant="body2" component="p">
            </Typography> */}
          </CardContent>
          <CardActions>
            <Button size="small" color="primary" href={`/meetings/${m.id}`}>
              join
            </Button>
          </CardActions>
          {/* <CardActions>
            <EditMessagePostForm
              prevMessage={item.message}
              id={item.id.toString()}
            />
            <DeleteMessageDialog
              targetMessage={item.message}
              id={item.id.toString()}
            />
          </CardActions> */}
        </Card>
      ))}
    </div>
  );
}

MeetingList.propTypes = {
  forceUpdate: PropTypes.number,
};

MeetingList.defaultProps = {
  forceUpdate: Math.random(),
};

export default function Meetings() {
  const [randomValue, setRandomValue] = useState<number>(Math.random());

  const onMeetingAdded = () => {
    setRandomValue(Math.random());
  };

  return (
    <>
      <CreateMeetingForm onSubmitSuccessful={onMeetingAdded} />
      <MeetingList forceUpdate={randomValue} />
    </>
  );
}
