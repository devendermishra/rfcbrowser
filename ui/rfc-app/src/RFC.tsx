import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { IRFCDetail } from './types';
import styles from './RFC.module.css'; // Import the styles
import DOMPurify from 'dompurify';

interface ParamTypes {
  id: string;
}

const RFC: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [rfc, setRfc] = useState<IRFCDetail | null>(null);
  
    useEffect(() => {
      fetch(`http://127.0.0.1:8080/rfc/${id}`)
        .then((response) => response.json())
        .then((data: IRFCDetail) => setRfc(data))
        .catch((error) => console.error('Error fetching RFC:', error));
    }, [id]);
  
    const createMarkup = (htmlContent: string) => {
        return { __html: DOMPurify.sanitize(htmlContent) };
    };
    const handleDownload = () => {
      fetch(`http://127.0.0.1:8080/download/rfc/${id}`, { method: 'POST' })
        .then(response => {
          if (response.ok) {
            console.log('RFC downloaded successfully');
            // Refresh the page to reflect any changes
            window.location.reload();
            // Optional: Update state or UI to reflect the download status
          } else {
            console.error('Failed to download RFC');
          }
        })
        .catch(error => console.error('Error downloading RFC:', error));
    };

    return (
      <div className={styles.container}>
        <button onClick={() => navigate(-1)} className={styles.backButton}>
          Back
        </button>
        <h1 className={styles.header}>{rfc?.id}</h1>
        <button onClick={handleDownload} className={styles.downloadButton}>
          Download
        </button>
        <div
        className={styles.content}
        dangerouslySetInnerHTML={createMarkup(rfc?.content || '')}></div>
        <button onClick={() => navigate(-1)} className={styles.backButton}>
          Back
        </button>
      </div>
    );
  };

export default RFC;
