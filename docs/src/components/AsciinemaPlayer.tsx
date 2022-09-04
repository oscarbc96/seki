import React, {FC, useEffect, useRef} from 'react';
import 'asciinema-player/dist/bundle/asciinema-player.css';
import BrowserOnly from "@docusaurus/BrowserOnly";

// Based on https://github.com/asciinema/asciinema-player/issues/72#issuecomment-1024929693

type AsciinemaPlayerProps = {
    src: string;
    // START asciinemaOptions
    cols: string;
    rows: string;
    autoPlay: boolean
    preload: boolean;
    loop: boolean | number;
    startAt: number | string;
    speed: number;
    idleTimeLimit: number;
    theme: string;
    poster: string;
    fit: boolean | string;
    terminalFontSize: string;
    // END asciinemaOptions
};


const AsciinemaPlayerBrowserOnly: FC<AsciinemaPlayerProps> = ({
    src,
    ...asciinemaOptions
}) => {
  return (
    <BrowserOnly fallback={<div>Loading...</div>}>
      {() => {
        const AsciinemaPlayerLibrary = require('asciinema-player').create;
        const ref = useRef<HTMLDivElement>(null);

        useEffect(() => {
            const currentRef = ref.current;
            AsciinemaPlayerLibrary(src, currentRef, asciinemaOptions);
        }, [src]);

        return <div ref={ref} />;
      }}
    </BrowserOnly>
  );
};

export default AsciinemaPlayerBrowserOnly;