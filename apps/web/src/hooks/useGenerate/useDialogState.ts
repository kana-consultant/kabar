export function closeResultDialog(
    setShowResultDialog: (val: boolean) => void,
    setPublishResults: (val: any) => void
) {
    setShowResultDialog(false);
    setPublishResults(null);
}