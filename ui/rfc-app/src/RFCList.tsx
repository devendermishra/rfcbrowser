import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { IRFC } from './types';
import styles from './RFCList.module.css'; // Import the styles

// Define the function that returns the RFC links
const renderRFCLinks = (rfcString: string | undefined, basePath: string) => {
    if (!rfcString) return null;
  
    return rfcString.split(',').map((rfcId, index, array) => (
      <React.Fragment key={rfcId.trim()}>
        <Link to={`${basePath}${rfcId.trim()}`}>
          {rfcId.trim()}
        </Link>
        {index < array.length - 1 ? ', ' : ''}
      </React.Fragment>
    ));
};

const RFCList: React.FC = () => {
  const [rfcs, setRfcs] = useState<IRFC[]>([]);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    fetch('http://127.0.0.1:8080/rfcs')
      .then((response) => response.json())
      .then((data: IRFC[]) => setRfcs(data))
      .catch((error) => console.error('Error fetching RFCs:', error));
  }, []);

  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(event.target.value.toLowerCase());
  };
  
  const filteredRFCs = rfcs.filter(rfc =>
    rfc.title.toLowerCase().includes(searchTerm) ||
    rfc.authors.toLowerCase().includes(searchTerm) ||
    rfc.id.toLowerCase().includes(searchTerm)
  );

  return (
    <div className={styles.container}>
      <h1 className={styles.header}>RFCs</h1>
      <input
        type="text"
        className={styles.searchBox}
        placeholder="Search by Title or Author..."
        onChange={handleSearchChange}
      />

      <ul className={styles.list}>
        <div key="rfc_id" className={styles.heading}>
            <div className={`${styles.cell} ${styles.cellId}`}>
                RFC Id
            </div>
            <div className={`${styles.cell} ${styles.cellTitle}`}>
                Title
            </div>
            <div className={`${styles.cell} ${styles.cellAuthors}`}>
                Authors
            </div>
            <div className={styles.cell}>
                Year/Month
            </div>
            <div className={styles.cell}>
                Status
            </div>
            <div className={styles.cell}>
                Obsoleted by
            </div>
            <div className={styles.cell}>
                Obsoletes
            </div>
            <div className={styles.cell}>
                Updated by
            </div>
            <div className={styles.cell}>
                Updates
            </div>
            <div className={styles.cell}>
                See Also
            </div>

        </div>
        {filteredRFCs.map((rfc) => (
            <div key={rfc.id} className={styles.row}>
                <div className={`${styles.cell} ${styles.cellId}`}>
                    {rfc.id}
                </div>
                <div className={`${styles.cell} ${styles.cellTitle}`}>
                    <Link to={`/rfc/${rfc.id}`} className={styles.link}>
                    {rfc.title}
                    </Link>
                </div>
                <div className={`${styles.cell} ${styles.cellAuthors}`}>
                    {rfc.authors}
                </div>
                <div className={styles.cell}>
                    {`${rfc.year}/${rfc.month}`}
                </div>
                <div className={styles.cell}>
                    {rfc.status}
                </div>
                <div className={styles.cell}>
                    {renderRFCLinks(rfc.obsoleted_by, '/rfc/')}
                </div>
                <div className={styles.cell}>
                    {renderRFCLinks(rfc.obsoletes, '/rfc/')}
                </div>
                <div className={styles.cell}>
                    {renderRFCLinks(rfc.updated_by, '/rfc/')}
                </div>
                <div className={styles.cell}>
                    {renderRFCLinks(rfc.updates, '/rfc/')}
                </div>
                <div className={styles.cell}>
                    {renderRFCLinks(rfc.also, '/rfc/')}
                </div>
          </div>
        ))}
      </ul>
    </div>
  );
};

export default RFCList;
