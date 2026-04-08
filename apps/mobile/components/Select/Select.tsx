import { AppSelect } from '@/components/ui/AppSelect';

interface SelectProps {
  onChange: (value: string) => void;
  data: { value: string; label: string }[];
  selected?: string | null;
}

export function Select({ onChange, data, selected }: SelectProps) {
  return <AppSelect options={data} value={selected} onValueChange={onChange} />;
}
