import React from 'react';

interface FormFieldProps {
    id: string;
    label: string;
    type: 'text' | 'email' | 'password';
    value: string;
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    placeholder?: string;
    required?: boolean;
}

export const FormField: React.FC<FormFieldProps> = ({
    id,
    label,
    type,
    value,
    onChange,
    placeholder,
    required = false,
}) => (
    <div className="form-group">
        <label htmlFor={id}>{label}</label>
        <input
            id={id}
            type={type}
            value={value}
            onChange={onChange}
            placeholder={placeholder}
            required={required}
        />
    </div>
);
