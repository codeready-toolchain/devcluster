import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import FormGroup from '@material-ui/core/FormGroup';
import FormControl from '@material-ui/core/FormControl';
import TextField from '@material-ui/core/TextField';
import Select from '@material-ui/core/Select';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import Button from '@material-ui/core/Button';

const useStyles = makeStyles((theme) => ({
  formControl: {
    margin: theme.spacing(1),
  },
  formRow: {
    marginTop: theme.spacing(1),
  },
}));

export default function RequestForm({ zones, onSubmit }) {
  const classes = useStyles();
  
  const [numberOfClusters, setNumberOfClusters] = React.useState(10);
  const [deleteInHours, setDeleteInHours] = React.useState(155);
  const [zone, setZone] = React.useState('');

  const onClickRequest = () => {
    onSubmit({
        numberOfClusters: numberOfClusters,
        zone: zone,
        deleteInHours: deleteInHours,
    })
  }

  React.useEffect(() => {
    if (zones && zones[0])
      setZone(zones[0].id);
  }, [zones]);

  return (
    <FormGroup className={classes.formRow} row>
        <FormControl className={classes.formControl} style={{minWidth: '220px'}}>
            <TextField
                value={numberOfClusters}
                label="Number of Clusters"
                type="number"
                InputProps={{
                  inputProps: { min: 1, max: 40 },
                }}
                InputLabelProps={{
                  shrink: true,
                }}
                onChange={(event) => event.target.value<1?setNumberOfClusters(1):event.target.value>40?setNumberOfClusters(40):setNumberOfClusters(event.target.value)}
            />
        </FormControl>
        <FormControl className={classes.formControl} style={{minWidth: '220px'}}>
            <TextField
                value={deleteInHours}
                label="Delete in Hours"
                type="number"
                InputProps={{
                  inputProps: { min: 1, max: 170 },
                }}
                InputLabelProps={{
                  shrink: true,
                }}
                onChange={(event) => event.target.value<1?setDeleteInHours(1):event.target.value>170?setDeleteInHours(170):setDeleteInHours(event.target.value)}
            />
        </FormControl>
        <FormControl className={classes.formControl} style={{minWidth: '220px'}}>
            <InputLabel id='zone-label'>Zone</InputLabel>
            <Select labelId='zone-label' id='zone-select' value={zone} onChange={(event) => setZone(event.target.value)}>
                {zones?zones.map((zone, index) =>
                    <MenuItem key={index} value={zone.id}>{zone['display_name']}</MenuItem>
                ):null}
            </Select>
        </FormControl>
        <FormControl className={classes.formControl}>
            <Button variant='contained' onClick={() => onClickRequest()}>Request Clusters</Button>
        </FormControl>
    </FormGroup>
  )
};
