import type { ReactNode } from 'react';

export function parseMessage(input: string): ReactNode {
    const parts = input.split('\n');
    return (
        <span className='whitespace-pre-line flex flex-col'>
            {parts.map((line, i) => (
                <span
                    key={i}
                    dangerouslySetInnerHTML={{ __html: line }}
                />
            ))}
        </span>
    );
}
