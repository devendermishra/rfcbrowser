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

  const refreshRFCs = () => {
    fetch('http://127.0.0.1:8080/refresh_rfcs', { method: 'POST' })
      .then((response) => response.json())
      .then((data: IRFC[]) => setRfcs(data))
      .catch((error) => console.error('Error fetching RFCs:', error));
  };

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
      <button onClick={refreshRFCs} className={styles.refreshButton}>
          Refresh List
      </button>
      <input
        type="text"
        className={styles.searchBox}
        placeholder="Search by Title or Author..."
        onChange={handleSearchChange}
      />
      <table className={styles.table}>
        <thead>
          <tr className={styles.heading}>
            <th className={styles.cellId}>RFC Id</th>
            <th className={styles.cellTitle}>Title</th>
            <th className={styles.cellAuthors}>Authors</th>
            <th className={styles.cell}>Year/Month</th>
            <th className={styles.cell}>Status</th>
            <th className={styles.cell}>Obsoleted by</th>
            <th className={styles.cell}>Obsoletes</th>
            <th className={styles.cell}>Updated by</th>
            <th className={styles.cell}>Updates</th>
            <th className={styles.cell}>See Also</th>
          </tr>
        </thead>
        <tbody>
          {filteredRFCs.map((rfc) => (
            <tr key={rfc.id} className={styles.row}>
              <td className={styles.cellId}>{rfc.id}</td>
              <td className={styles.cellTitle}>
                <Link to={`/rfc/${rfc.id}`} className={styles.link}>
                  {rfc.title}
                </Link>
              </td>
              <td className={styles.cellAuthors}>{rfc.authors}</td>
              <td className={styles.cell}>{`${rfc.year}/${rfc.month}`}</td>
              <td className={styles.cell}>{rfc.status}</td>
              <td className={styles.cell}>
                {renderRFCLinks(rfc.obsoleted_by, '/rfc/')}
              </td>
              <td className={styles.cell}>
                {renderRFCLinks(rfc.obsoletes, '/rfc/')}
              </td>
              <td className={styles.cell}>
                {renderRFCLinks(rfc.updated_by, '/rfc/')}
              </td>
              <td className={styles.cell}>
                {renderRFCLinks(rfc.updates, '/rfc/')}
              </td>
              <td className={styles.cell}>
                {renderRFCLinks(rfc.also, '/rfc/')}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default RFCList;
