// app/loading.tsx
// Глобальный индикатор загрузки для приложения (серверный компонент).

import Spinner from "../utils/spinner/Spinner";
import styles from "../utils/spinner/spinner.module.css";

export default function Loading() {
    return (
        <div className={styles.spinnerContainer}>
            <Spinner />
        </div>
    );
}