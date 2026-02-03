import { TextareaHTMLAttributes } from "preact";

type TextareaProps = TextareaHTMLAttributes & {
    error?: string;
};

const base = `inline-flex w-full min-h-[96px] px-4 py-2 text-sm mt-1 font-medium border bg-zinc-900 transition outline-none resize-y disabled:pointer-events-none disabled:opacity-40`;
const primary = `placeholder:text-zinc-500   text-zinc-100   border-zinc-700   focus:border-zinc-300   focus:text-zinc-300`;
const error_ = ` placeholder:text-orange-800 text-orange-600 border-orange-600 focus:border-orange-600 focus:text-orange-600`;

export function Textarea({ class: className, error, ...props }: TextareaProps) {
    return (
        <div class="w-full">
            <textarea
                {...props}
                class={[base, error ? error_ : primary, className].filter(Boolean).join(" ")}
            />

            {error && <p class="mt-1 text-sm text-orange-600">{error}</p>}
        </div>
    );
}
