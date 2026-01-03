/* eslint-disable react-hooks/set-state-in-effect */
import { useMutation } from '@tanstack/react-query';
import { scanFile, type ScanResult } from '../server/api';
import DropBox from '../components/dropbox';
import { illustrations } from '../components/illustrations';
import { useEffect, useState, type ReactNode } from 'react';
import { highlightText } from '../utils/highlightText';
import { parseMessage } from '../utils/parseMessage';

export default function App() {
    const [file, setFile] = useState<File | null>(null);
    const [message, setMessage] = useState<ReactNode>(<></>);
    const mutation = useMutation<ScanResult, Error, File>({
        mutationFn: scanFile,
    });

    useEffect(() => {
        if (mutation.isPending) {
            setMessage(
                parseMessage(
                    `${highlightText(
                        file?.name || 'Your File',
                        'info'
                    )} is currently being scanned, please wait a moment`
                )
            );
            return;
        }

        if (mutation.isError) {
            console.log(JSON.stringify(mutation.error));
            setMessage(
                <>
                    <ResultView
                        result={mutation.error}
                        fileName={file?.name || 'Your file'}
                    />
                </>
            );
            return;
        }

        if (mutation.data) {
            console.log(mutation);
            setMessage(
                <ResultView
                    result={mutation.data}
                    fileName={file?.name || 'Your file'}
                />
            );
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [
        file?.name,
        mutation.isPending,
        mutation.isError,
        mutation.data,
        mutation.error,
    ]);

    return (
        <div className='w-full h-full flex lg:items-center justify-center bg-base-300 overflow-hidden text-base-content'>
            <main className='relative top-8 mb-8 flex h-fit w-11/12 lg:w-6/12 items-center justify-center rounded-lg lg:px-8 py-4 lg:py-8 shadow 2xl:w-50% 2xl:min-w-50% bg-base-100'>
                <div className='flex w-11/12 flex-col gap-4 items-center justify-center'>
                    <h1 className='text-2xl font-semibold mb-4 items-center flex'>
                        Virus Scanner
                    </h1>
                    <DropBox
                        fileType='*/*'
                        illustration={
                            <illustrations.Secure
                                isError={
                                    mutation.data?.status.toLowerCase() ===
                                    'infected'
                                }
                            />
                        }
                        callback={(file) => {
                            setFile(file);
                            mutation.mutate(file);
                        }}
                        contentMessage={message}
                    />
                </div>
            </main>
        </div>
    );
}

function ResultView({
    result,
    fileName,
}: {
    result: ScanResult | string;
    fileName: string;
}) {
    if (typeof result === 'string') {
        return (
            <p className='mt-4 text-error whitespace-pre-line break-words text-pretty w-full'>
                {result ? `ERROR: ${result}` : 'Error during scan'}
            </p>
        );
    }

    if (result.status === 'clean') {
        return (
            <p className='mt-4 text-success font-medium whitespace-pre-line break-words text-pretty w-full'>
                {fileName} is clean
            </p>
        );
    }

    if (result.status === 'infected') {
        return (
            <p className='mt-4 text-error font-medium whitespace-pre-line break-words text-pretty w-full'>
                {fileName} is infected ({result.signature})
            </p>
        );
    }
}
