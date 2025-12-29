/* eslint-disable react-hooks/set-state-in-effect */
import {
    useEffect,
    useRef,
    useState,
    type ReactNode,
    type RefObject,
} from 'react';
import { cn } from '../utils/cn';
import { highlightText } from '../utils/highlightText';
import { useIsMobile } from '../hooks/useIsMobile';
import { parseMessage } from '../utils/parseMessage';

export type MimeType = string;

type Props = {
    fileType: MimeType;
    illustration: ReactNode;
    url?: string | null;
    callback: (file: File) => void;
    willLoad?: boolean;
    loading?: boolean;
    progress?: number;
    contentMessage?: ReactNode;
};

export default function DropBox({
    fileType,
    illustration,
    url,
    callback,
    willLoad,
    loading,
    progress,
    contentMessage,
}: Props) {
    const isMobile = useIsMobile();
    const [message, setMessage] = useState<ReactNode>('');
    const [fileName, setFileName] = useState<string | undefined>();
    const fileInputRef = useRef<HTMLInputElement>(null);
    const hiddenDropAreaRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (contentMessage) {
            setMessage(contentMessage);
        }
    }, [contentMessage]);

    useEffect(() => {
        const base = isMobile
            ? parseMessage(
                  `${highlightText(
                      'Click here'
                  )} to browse\n for a document to upload`
              )
            : parseMessage(
                  `Drag and drop your file here\n\nOr ${highlightText(
                      'click here'
                  )} to browse\nyour PC for documents to upload`
              );
        setMessage(base);
    }, [isMobile]);

    useEffect(() => {
        if (loading && typeof progress === 'number' && progress > 0) {
            const msg = `${highlightText(
                fileName || 'Your File',
                'info'
            )} is currently being scanned, please wait a moment\n\n${progress.toFixed(
                2
            )}% Complete`;
            setMessage(parseMessage(msg));
        }
        if (typeof progress === 'number' && progress > 99) {
            const msg = `${highlightText(
                fileName || 'Your File',
                'success'
            )} has finished scanning`;
            setMessage(parseMessage(msg));
        }
    }, [loading, progress, fileName]);

    function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
        const file = e.target.files?.[0];
        if (!file) {
            return;
        }

        const validity = checkFileTypeValidity(fileType, file);
        if (validity.valid) {
            setFileName(file.name);
            const parsed = parseMessage(validity.message);
            if (!willLoad) {
                setMessage(parsed);
            }
            callback(file);
        } else {
            setMessage(parseMessage(validity.message));
        }
    }

    return (
        <>
            <div className='relative z-10 mx-auto flex h-fit w-full flex-row items-center rounded-lg bg-base-200'>
                <div className='flex h-full w-full flex-col-reverse items-center justify-start pt-4 md:flex-row md:pt-0'>
                    <div
                        className={cn(
                            'relative flex h-full w-full rounded-lg bg-inherit md:w-[70%]'
                        )}>
                        <div
                            className={cn(
                                'm-5 h-full w-full cursor-pointer rounded-lg bg-base-100 p-4 lg:p-[60px_80px]',
                                url ? '!pb-[143px]' : ''
                            )}
                            onClick={() => {
                                if (!url) {
                                    fileInputRef.current?.click();
                                }
                            }}>
                            {illustration}
                        </div>
                    </div>

                    <label
                        className='font-size-1-25 flex h-full w-full cursor-pointer flex-col items-center justify-center text-pretty break-words pr-1 text-left md:w-[27%]'
                        onClick={() => {
                            fileInputRef.current?.click();
                        }}>
                        {message}
                        <input
                            type='file'
                            className='hidden'
                            ref={fileInputRef}
                            accept={fileType}
                            onChange={handleFileChange}
                        />
                    </label>
                </div>
            </div>

            <div
                className='absolute left-0 top-0 z-0 inline-block h-full w-full cursor-default rounded-lg opacity-0 outline-dashed outline-4 outline-accent focus:opacity-50 focus-visible:z-40'
                tabIndex={0}
                role='button'
                ref={hiddenDropAreaRef}
                data-id-name='hiddenDropZone'
                onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                        fileInputRef.current?.click();
                    }
                }}
                onMouseDown={(e) => e.preventDefault()}
                onDragOver={(e) =>
                    uploadBoxDropOverOrEnter(
                        e,
                        hiddenDropAreaRef as RefObject<HTMLDivElement>
                    )
                }
                onDragEnter={(e) =>
                    uploadBoxDropOverOrEnter(
                        e,
                        hiddenDropAreaRef as RefObject<HTMLDivElement>
                    )
                }
                onDragLeave={(e) =>
                    removeDropZone(
                        e,
                        hiddenDropAreaRef as RefObject<HTMLDivElement>
                    )
                }
                onDrop={(e) =>
                    uploadBoxOnDrop(
                        e,
                        hiddenDropAreaRef as RefObject<HTMLDivElement>,
                        fileInputRef as RefObject<HTMLInputElement>
                    )
                }
            />
        </>
    );
}

function checkFileTypeValidity(
    accept: string,
    file: File
): { valid: boolean; message: string } {
    if (!accept || accept === '*/*') {
        return {
            valid: true,
            message: `${highlightText(
                file.name,
                'success'
            )} is ready to upload`,
        };
    }

    const allowed = accept
        .split(',')
        .map((s) => s.trim())
        .filter((s) => s.length > 0);

    const ok = allowed.some((rule) => matchesAcceptRule(rule, file));
    if (ok) {
        return {
            valid: true,
            message: `${highlightText(
                file.name,
                'success'
            )} is ready to upload`,
        };
    }

    return {
        valid: false,
        message: `File type not allowed\n\nAllowed: ${allowed.join(', ')}`,
    };
}

function matchesAcceptRule(rule: string, file: File) {
    if (rule.endsWith('/*')) {
        const prefix = rule.slice(0, -1);
        return file.type.startsWith(prefix);
    }
    if (rule.startsWith('.')) {
        return file.name.toLowerCase().endsWith(rule.toLowerCase());
    }
    return file.type === rule;
}

function uploadBoxDropOverOrEnter(
    e: React.DragEvent,
    areaRef: React.RefObject<HTMLDivElement>
) {
    e.preventDefault();
    const el = areaRef.current;
    if (el) {
        el.classList.add('opacity-50');
    }
}

function removeDropZone(
    e: React.DragEvent,
    areaRef: React.RefObject<HTMLDivElement>
) {
    e.preventDefault();
    const el = areaRef.current;
    if (el) {
        el.classList.remove('opacity-50');
    }
}

function uploadBoxOnDrop(
    e: React.DragEvent,
    areaRef: React.RefObject<HTMLDivElement>,
    inputRef: React.RefObject<HTMLInputElement>
) {
    e.preventDefault();
    const el = areaRef.current;
    if (el) {
        el.classList.remove('opacity-50');
    }
    const file = e.dataTransfer.files?.[0];
    if (file && inputRef.current) {
        const dt = new DataTransfer();
        dt.items.add(file);
        inputRef.current.files = dt.files;
        inputRef.current.dispatchEvent(new Event('change', { bubbles: true }));
    }
}
