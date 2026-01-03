export type ScanResult = {
    status: 'clean' | 'infected' | 'error';
    signature?: string;
    engine: string;
    duration_ms: number;
    error?: string;
};

export async function scanFile(file: File): Promise<ScanResult> {
    const form = new FormData();
    form.append('file', file);
    // form.append(
    //     'content',
    //     'X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*'
    // );

    const res = await fetch('/api/scan', {
        method: 'POST',
        body: form,
    });

    const body = await res.json();

    if (!res.ok) {
        throw body.error;
    }

    return body.message;
}

export async function apiHealthCheck() {
    const res = await fetch('/api/health');

    if (!res.ok) {
        const body = await res.json();
        throw new Error(body?.error);
    }

    const body = await res.json();
    return body.message;
}
