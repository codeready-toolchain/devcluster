import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import IconButton from '@material-ui/core/IconButton';
import CloudDownloadIcon from '@material-ui/icons/CloudDownload';
import FormGroup from '@material-ui/core/FormGroup';
import FormControl from '@material-ui/core/FormControl';
import Slider from '@material-ui/core/Slider';
import InputLabel from '@material-ui/core/InputLabel';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';

const useStyles = makeStyles((theme) => ({
  container: {
    display: 'flex',
    flex: 1,
    flexDirection: 'row',
  },
  formControl: {
    margin: theme.spacing(1),
  },
  form: {
    display: 'flex',
    flex: 1,
    flexDirection: 'row',
    marginTop: theme.spacing(1),
  },
  exportButton: {
    display: 'flex',
    flex: 1,
    justifyContent: 'flex-end',
  }
}));

export default function UserForm({ onSubmit, onExport }) {
  const classes = useStyles();
  
  const [numberOfUsers, setNumberOfUsers] = React.useState(10);
  const [startIndex, setStartIndex] = React.useState(0);

  const onClickRequest = () => {
    onSubmit({
        numberOfUsers: numberOfUsers,
        startIndex: startIndex,
    })
  }

  return (
    <div className={classes.container}>
      <div className={classes.form}>
        <FormGroup row>
            <FormControl className={classes.formControl} style={{minWidth: '220px'}}>
                <InputLabel id='usersn-label'>Number of Users</InputLabel>
                <Slider
                    defaultValue={10}
                    aria-labelledby='usern-label'
                    valueLabelDisplay='auto'
                    step={10}
                    marks
                    min={10}
                    max={500}
                    onChange={(event, newValue) => setNumberOfUsers(newValue)}
                />
            </FormControl>
            <FormControl className={classes.formControl} style={{minWidth: '220px'}}>
                <TextField
                    label="Start Index"
                    defaultValue={0}
                    type="number"
                    InputLabelProps={{
                      shrink: true,
                    }}
                    onChange={(event) => setStartIndex(event.target.value)}
                />
            </FormControl>
            <FormControl className={classes.formControl}>
                <Button variant='contained' onClick={() => onClickRequest()}>Request Users</Button>
            </FormControl>
        </FormGroup>
      </div>
      <div className={classes.exportButton} >
        <IconButton aria-label='export' color='primary' onClick={() => onExport()}>
            <CloudDownloadIcon/>
        </IconButton>
      </div>
    </div>
  )
};
