export interface Field {
    id: string;
    key: string;
    value: string;
    type: "field" | "object";
    children: Field[];
    expanded: boolean;
}

export interface SimpleJsonBuilderProps {
    value: any;
    onChange: (value: any) => void;
}