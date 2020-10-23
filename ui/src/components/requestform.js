import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import FormGroup from '@material-ui/core/FormGroup';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Slider from '@material-ui/core/Slider';
import Select from '@material-ui/core/Select';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import Checkbox from '@material-ui/core/Checkbox';
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
  const [deleteInHours, setDeleteInHours] = React.useState(24);
  const [zone, setZone] = React.useState('');
  const [subnet, setSubnet] = React.useState(false);

  const onClickRequest = () => {
    onSubmit({
        numberOfClusters: numberOfClusters,
        zone: zone,
        deleteInHours: deleteInHours,
        subnet: subnet,
    })
  }

  return (
    <FormGroup className={classes.formRow} row>
        <FormControl className={classes.formControl} style={{minWidth: '220px'}}>
            <InputLabel id='clustern-label'>Number of Clusters</InputLabel>
            <Slider
                defaultValue={1}
                aria-labelledby='clustern-label'
                valueLabelDisplay='auto'
                step={1}
                marks
                min={1}
                max={200}
                onChange={(event, newValue) => setNumberOfClusters(newValue)}
            />
        </FormControl>
        <FormControl className={classes.formControl} style={{minWidth: '220px'}}>
            <InputLabel id='ttl-label'>Delete in Hours</InputLabel>
            <Slider
                defaultValue={10}
                aria-labelledby='ttl-label'
                valueLabelDisplay='auto'
                step={24}
                marks
                min={24}
                max={720}
                onChange={(event, newValue) => setDeleteInHours(newValue)}
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