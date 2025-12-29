import { useEffect, useState } from 'react';

export function useIsMobile(breakpoint = 768) {
    const [isMobile, setIsMobile] = useState(true);

    useEffect(() => {
        function checkIsMobile() {
            setIsMobile(window.innerWidth < breakpoint);
        }
        checkIsMobile();

        window.addEventListener('resize', checkIsMobile);
        return () => {
            window.removeEventListener('resize', checkIsMobile);
        };
    }, [breakpoint]);

    return isMobile;
}
