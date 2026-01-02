import type { ReactNode } from 'react';

export function parseMessage(input: string): ReactNode {
    const parts = input.split('\n');
    return (
        <span className='whitespace-pre-line flex flex-col breaks-words text-wrap w-full'>
            {parts.map((line, i) => (
                <span
                    key={i}
                    className='whitespace-pre-line break-words text-pretty w-full'
                    dangerouslySetInnerHTML={{ __html: line }}
                />
            ))}
        </span>
    );
}
