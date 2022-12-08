import React, { useEffect } from 'react';

const CustomHTMLPost = (props: any) => {
    console.log(props);

    return (
        <div
            className='custom_post_container'
            dangerouslySetInnerHTML={{__html: props.post.message_source}}
        />
    )
}

export default CustomHTMLPost;
