import React, {Component} from 'react';
import PropTypes from 'prop-types';

class AlbumItem extends Component {
    render() {
        return (
            <div>{this.props.album.title}</div>
        )
    }
}

AlbumItem.propTypes = {
    album: PropTypes.object.isRequired
};

export default AlbumItem;
