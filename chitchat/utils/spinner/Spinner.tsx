//  Spinner - вращающаяся иконка загрузки.

import Image from 'next/image';
import styles from './spinner.module.css';

interface SpinnerProps {
  size?: 'small' | 'large';
  className?: string;
}

export default function Spinner({ size = 'small', className }: SpinnerProps) {
  const spinnerClass = size === 'large' ? styles.spinnerLarge : styles.spinner;
  
  return (
    <Image
      src="/spinner.svg"
      alt="Loading..."
      width={size === 'large' ? 48 : 24}
      height={size === 'large' ? 48 : 24}
      className={`${spinnerClass} ${className || ''}`}
    />
  );
}

export function SpinnerContainer({ size = 'large' }: { size?: 'small' | 'large' }) {
  return (
    <div className={styles.spinnerContainer}>
      <Spinner size={size} />
    </div>
  );
}