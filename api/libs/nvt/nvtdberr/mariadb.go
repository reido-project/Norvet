package nvtdberr

import (
	"api/libs/nvt/nvterr"

	"github.com/go-sql-driver/mysql"
)

// identifyMariaDB translates a MariaDB error number into an nvterr error.
// Unmapped numbers return nvterr.ErrInternal.
func identifyMariaDB(me *mysql.MySQLError) (*nvterr.Error, Info) {
	return identifyMySQLFamily(mariadbErrors, me)
}

// mariadbErrors maps MariaDB server error numbers (1000-4999) and MariaDB
// client-library error numbers (5000+) to client-safe nvterr errors.
// Only 4xx and non-500 5xx errors are listed: every unlisted number falls back
// to nvterr.ErrInternal. Entries up to 1969 are annotated with the nvterr code
// and status; later entries with the MariaDB symbolic name.
var mariadbErrors = map[uint16]*nvterr.Error{
	1007: nvterr.ErrAlreadyExists,       // already_exists (409)
	1008: nvterr.ErrNotFound,            // not_found (404)
	1022: nvterr.ErrDuplicateValue,      // duplicate_value (422)
	1046: nvterr.ErrBadRequest,          // bad_request (400)
	1048: nvterr.ErrRequired,            // required (422)
	1049: nvterr.ErrNotFound,            // not_found (404)
	1050: nvterr.ErrAlreadyExists,       // already_exists (409)
	1051: nvterr.ErrNotFound,            // not_found (404)
	1052: nvterr.ErrBadRequest,          // bad_request (400)
	1054: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1055: nvterr.ErrBadRequest,          // bad_request (400)
	1056: nvterr.ErrBadRequest,          // bad_request (400)
	1057: nvterr.ErrBadRequest,          // bad_request (400)
	1058: nvterr.ErrBadRequest,          // bad_request (400)
	1059: nvterr.ErrValueTooLong,        // value_too_long (422)
	1060: nvterr.ErrDuplicateValue,      // duplicate_value (422)
	1061: nvterr.ErrAlreadyExists,       // already_exists (409)
	1062: nvterr.ErrDuplicateValue,      // duplicate_value (422)
	1063: nvterr.ErrBadRequest,          // bad_request (400)
	1064: nvterr.ErrBadRequest,          // bad_request (400)
	1065: nvterr.ErrBadRequest,          // bad_request (400)
	1066: nvterr.ErrBadRequest,          // bad_request (400)
	1067: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1068: nvterr.ErrBadRequest,          // bad_request (400)
	1069: nvterr.ErrBadRequest,          // bad_request (400)
	1070: nvterr.ErrBadRequest,          // bad_request (400)
	1071: nvterr.ErrValueTooLong,        // value_too_long (422)
	1072: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1073: nvterr.ErrBadRequest,          // bad_request (400)
	1074: nvterr.ErrValueTooLarge,       // value_too_large (422)
	1075: nvterr.ErrBadRequest,          // bad_request (400)
	1082: nvterr.ErrNotFound,            // not_found (404)
	1083: nvterr.ErrBadRequest,          // bad_request (400)
	1084: nvterr.ErrBadRequest,          // bad_request (400)
	1086: nvterr.ErrAlreadyExists,       // already_exists (409)
	1089: nvterr.ErrBadRequest,          // bad_request (400)
	1090: nvterr.ErrBadRequest,          // bad_request (400)
	1091: nvterr.ErrNotFound,            // not_found (404)
	1093: nvterr.ErrBadRequest,          // bad_request (400)
	1094: nvterr.ErrNotFound,            // not_found (404)
	1096: nvterr.ErrBadRequest,          // bad_request (400)
	1097: nvterr.ErrValueTooLarge,       // value_too_large (422)
	1099: nvterr.ErrLocked,              // locked (423)
	1102: nvterr.ErrBadRequest,          // bad_request (400)
	1103: nvterr.ErrBadRequest,          // bad_request (400)
	1106: nvterr.ErrNotFound,            // not_found (404)
	1107: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1108: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1109: nvterr.ErrNotFound,            // not_found (404)
	1110: nvterr.ErrBadRequest,          // bad_request (400)
	1111: nvterr.ErrBadRequest,          // bad_request (400)
	1112: nvterr.ErrNotSupported,        // not_supported (422)
	1113: nvterr.ErrBadRequest,          // bad_request (400)
	1115: nvterr.ErrNotSupported,        // not_supported (422)
	1118: nvterr.ErrValueTooLarge,       // value_too_large (422)
	1120: nvterr.ErrBadRequest,          // bad_request (400)
	1121: nvterr.ErrNotSupported,        // not_supported (422)
	1125: nvterr.ErrAlreadyExists,       // already_exists (409)
	1128: nvterr.ErrNotFound,            // not_found (404)
	1136: nvterr.ErrBadRequest,          // bad_request (400)
	1138: nvterr.ErrBadRequest,          // bad_request (400)
	1139: nvterr.ErrInvalidFormat,       // invalid_format (422)
	1140: nvterr.ErrBadRequest,          // bad_request (400)
	1146: nvterr.ErrNotFound,            // not_found (404)
	1149: nvterr.ErrBadRequest,          // bad_request (400)
	1162: nvterr.ErrValueTooLong,        // value_too_long (422)
	1163: nvterr.ErrNotSupported,        // not_supported (422)
	1164: nvterr.ErrNotSupported,        // not_supported (422)
	1166: nvterr.ErrBadRequest,          // bad_request (400)
	1167: nvterr.ErrBadRequest,          // bad_request (400)
	1169: nvterr.ErrDuplicateValue,      // duplicate_value (422)
	1170: nvterr.ErrBadRequest,          // bad_request (400)
	1171: nvterr.ErrBadRequest,          // bad_request (400)
	1172: nvterr.ErrBadRequest,          // bad_request (400)
	1173: nvterr.ErrBadRequest,          // bad_request (400)
	1175: nvterr.ErrPreconditionFailed,  // precondition_failed (412)
	1176: nvterr.ErrNotFound,            // not_found (404)
	1177: nvterr.ErrNotFound,            // not_found (404)
	1178: nvterr.ErrNotImplemented,      // not_implemented (501)
	1179: nvterr.ErrInvalidState,        // invalid_state (409)
	1191: nvterr.ErrNotFound,            // not_found (404)
	1192: nvterr.ErrInvalidState,        // invalid_state (409)
	1193: nvterr.ErrNotFound,            // not_found (404)
	1210: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1215: nvterr.ErrConflict,            // conflict (409)
	1216: nvterr.ErrConflict,            // conflict (409)
	1217: nvterr.ErrConflict,            // conflict (409)
	1221: nvterr.ErrBadRequest,          // bad_request (400)
	1222: nvterr.ErrBadRequest,          // bad_request (400)
	1224: nvterr.ErrBadRequest,          // bad_request (400)
	1225: nvterr.ErrDuplicateValue,      // duplicate_value (422)
	1230: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1231: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1232: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1234: nvterr.ErrBadRequest,          // bad_request (400)
	1235: nvterr.ErrNotSupported,        // not_supported (422)
	1238: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1239: nvterr.ErrBadRequest,          // bad_request (400)
	1240: nvterr.ErrBadRequest,          // bad_request (400)
	1241: nvterr.ErrBadRequest,          // bad_request (400)
	1242: nvterr.ErrBadRequest,          // bad_request (400)
	1247: nvterr.ErrBadRequest,          // bad_request (400)
	1248: nvterr.ErrBadRequest,          // bad_request (400)
	1250: nvterr.ErrBadRequest,          // bad_request (400)
	1252: nvterr.ErrBadRequest,          // bad_request (400)
	1253: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1256: nvterr.ErrValueTooLarge,       // value_too_large (422)
	1263: nvterr.ErrRequired,            // required (422)
	1264: nvterr.ErrValueOutOfRange,     // value_out_of_range (422)
	1265: nvterr.ErrValueTooLong,        // value_too_long (422)
	1267: nvterr.ErrBadRequest,          // bad_request (400)
	1270: nvterr.ErrBadRequest,          // bad_request (400)
	1271: nvterr.ErrBadRequest,          // bad_request (400)
	1273: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1280: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1281: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1283: nvterr.ErrBadRequest,          // bad_request (400)
	1286: nvterr.ErrNotSupported,        // not_supported (422)
	1288: nvterr.ErrBadRequest,          // bad_request (400)
	1289: nvterr.ErrNotSupported,        // not_supported (422)
	1291: nvterr.ErrDuplicateValue,      // duplicate_value (422)
	1292: nvterr.ErrInvalidDate,         // invalid_date (422)
	1295: nvterr.ErrNotSupported,        // not_supported (422)
	1298: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1299: nvterr.ErrInvalidDate,         // invalid_date (422)
	1304: nvterr.ErrAlreadyExists,       // already_exists (409)
	1305: nvterr.ErrNotFound,            // not_found (404)
	1312: nvterr.ErrNotSupported,        // not_supported (422)
	1314: nvterr.ErrNotSupported,        // not_supported (422)
	1335: nvterr.ErrNotImplemented,      // not_implemented (501)
	1336: nvterr.ErrNotSupported,        // not_supported (422)
	1348: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1359: nvterr.ErrAlreadyExists,       // already_exists (409)
	1360: nvterr.ErrNotFound,            // not_found (404)
	1364: nvterr.ErrRequired,            // required (422)
	1365: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1366: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1367: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1369: nvterr.ErrPreconditionFailed,  // precondition_failed (412)
	1372: nvterr.ErrInvalidLength,       // invalid_length (422)
	1406: nvterr.ErrValueTooLong,        // value_too_long (422)
	1411: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1416: nvterr.ErrInvalidFormat,       // invalid_format (422)
	1425: nvterr.ErrValueTooLarge,       // value_too_large (422)
	1426: nvterr.ErrValueTooLarge,       // value_too_large (422)
	1427: nvterr.ErrBadRequest,          // bad_request (400)
	1431: nvterr.ErrNotFound,            // not_found (404)
	1432: nvterr.ErrInvalidFormat,       // invalid_format (422)
	1433: nvterr.ErrInvalidFormat,       // invalid_format (422)
	1437: nvterr.ErrValueTooLong,        // value_too_long (422)
	1439: nvterr.ErrValueOutOfRange,     // value_out_of_range (422)
	1441: nvterr.ErrInvalidDate,         // invalid_date (422)
	1446: nvterr.ErrInvalidFormat,       // invalid_format (422)
	1451: nvterr.ErrConflict,            // conflict (409)
	1452: nvterr.ErrConflict,            // conflict (409)
	1458: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1463: nvterr.ErrBadRequest,          // bad_request (400)
	1468: nvterr.ErrValueTooLong,        // value_too_long (422)
	1469: nvterr.ErrValueTooLong,        // value_too_long (422)
	1470: nvterr.ErrInvalidLength,       // invalid_length (422)
	1474: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1475: nvterr.ErrBadRequest,          // bad_request (400)
	1476: nvterr.ErrAlreadyExists,       // already_exists (409)
	1477: nvterr.ErrNotFound,            // not_found (404)
	1478: nvterr.ErrNotSupported,        // not_supported (422)
	1479: nvterr.ErrBadRequest,          // bad_request (400)
	1480: nvterr.ErrBadRequest,          // bad_request (400)
	1481: nvterr.ErrBadRequest,          // bad_request (400)
	1482: nvterr.ErrBadRequest,          // bad_request (400)
	1483: nvterr.ErrBadRequest,          // bad_request (400)
	1484: nvterr.ErrBadRequest,          // bad_request (400)
	1485: nvterr.ErrBadRequest,          // bad_request (400)
	1486: nvterr.ErrBadRequest,          // bad_request (400)
	1487: nvterr.ErrBadRequest,          // bad_request (400)
	1488: nvterr.ErrNotFound,            // not_found (404)
	1489: nvterr.ErrBadRequest,          // bad_request (400)
	1490: nvterr.ErrBadRequest,          // bad_request (400)
	1491: nvterr.ErrNotSupported,        // not_supported (422)
	1492: nvterr.ErrBadRequest,          // bad_request (400)
	1493: nvterr.ErrBadRequest,          // bad_request (400)
	1494: nvterr.ErrBadRequest,          // bad_request (400)
	1495: nvterr.ErrDuplicateValue,      // duplicate_value (422)
	1496: nvterr.ErrBadRequest,          // bad_request (400)
	1498: nvterr.ErrNotFound,            // not_found (404)
	1499: nvterr.ErrValueTooLarge,       // value_too_large (422)
	1500: nvterr.ErrBadRequest,          // bad_request (400)
	1502: nvterr.ErrBadRequest,          // bad_request (400)
	1503: nvterr.ErrBadRequest,          // bad_request (400)
	1504: nvterr.ErrBadRequest,          // bad_request (400)
	1505: nvterr.ErrBadRequest,          // bad_request (400)
	1506: nvterr.ErrNotSupported,        // not_supported (422)
	1507: nvterr.ErrNotFound,            // not_found (404)
	1508: nvterr.ErrBadRequest,          // bad_request (400)
	1509: nvterr.ErrBadRequest,          // bad_request (400)
	1510: nvterr.ErrBadRequest,          // bad_request (400)
	1511: nvterr.ErrBadRequest,          // bad_request (400)
	1512: nvterr.ErrBadRequest,          // bad_request (400)
	1513: nvterr.ErrBadRequest,          // bad_request (400)
	1514: nvterr.ErrBadRequest,          // bad_request (400)
	1515: nvterr.ErrBadRequest,          // bad_request (400)
	1516: nvterr.ErrNotFound,            // not_found (404)
	1517: nvterr.ErrAlreadyExists,       // already_exists (409)
	1519: nvterr.ErrBadRequest,          // bad_request (400)
	1520: nvterr.ErrBadRequest,          // bad_request (400)
	1523: nvterr.ErrBadRequest,          // bad_request (400)
	1525: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1526: nvterr.ErrValueOutOfRange,     // value_out_of_range (422)
	1527: nvterr.ErrBadRequest,          // bad_request (400)
	1531: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1532: nvterr.ErrValueTooLarge,       // value_too_large (422)
	1537: nvterr.ErrAlreadyExists,       // already_exists (409)
	1539: nvterr.ErrNotFound,            // not_found (404)
	1542: nvterr.ErrValueOutOfRange,     // value_out_of_range (422)
	1543: nvterr.ErrBadRequest,          // bad_request (400)
	1544: nvterr.ErrBadRequest,          // bad_request (400)
	1546: nvterr.ErrBadRequest,          // bad_request (400)
	1550: nvterr.ErrBadRequest,          // bad_request (400)
	1551: nvterr.ErrAlreadyExists,       // already_exists (409)
	1552: nvterr.ErrValueTooLong,        // value_too_long (422)
	1553: nvterr.ErrConflict,            // conflict (409)
	1554: nvterr.ErrBadRequest,          // bad_request (400)
	1562: nvterr.ErrBadRequest,          // bad_request (400)
	1563: nvterr.ErrBadRequest,          // bad_request (400)
	1564: nvterr.ErrBadRequest,          // bad_request (400)
	1566: nvterr.ErrBadRequest,          // bad_request (400)
	1567: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1568: nvterr.ErrInvalidState,        // invalid_state (409)
	1569: nvterr.ErrDuplicateValue,      // duplicate_value (422)
	1575: nvterr.ErrInvalidFormat,       // invalid_format (422)
	1576: nvterr.ErrBadRequest,          // bad_request (400)
	1578: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1579: nvterr.ErrNotSupported,        // not_supported (422)
	1580: nvterr.ErrBadRequest,          // bad_request (400)
	1581: nvterr.ErrBadRequest,          // bad_request (400)
	1582: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1583: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1584: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1585: nvterr.ErrConflict,            // conflict (409)
	1586: nvterr.ErrDuplicateValue,      // duplicate_value (422)
	1588: nvterr.ErrBadRequest,          // bad_request (400)
	1589: nvterr.ErrBadRequest,          // bad_request (400)
	1591: nvterr.ErrValueOutOfRange,     // value_out_of_range (422)
	1611: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1616: nvterr.ErrNotSupported,        // not_supported (422)
	1628: nvterr.ErrValueTooLong,        // value_too_long (422)
	1629: nvterr.ErrValueTooLong,        // value_too_long (422)
	1630: nvterr.ErrBadRequest,          // bad_request (400)
	1631: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1632: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1633: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1634: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1635: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1636: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1638: nvterr.ErrNotImplemented,      // not_implemented (501)
	1641: nvterr.ErrBadRequest,          // bad_request (400)
	1643: nvterr.ErrNotFound,            // not_found (404)
	1645: nvterr.ErrBadRequest,          // bad_request (400)
	1646: nvterr.ErrBadRequest,          // bad_request (400)
	1647: nvterr.ErrValueTooLong,        // value_too_long (422)
	1648: nvterr.ErrValueTooLong,        // value_too_long (422)
	1649: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1652: nvterr.ErrDuplicateValue,      // duplicate_value (422)
	1653: nvterr.ErrBadRequest,          // bad_request (400)
	1654: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1655: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1656: nvterr.ErrBadRequest,          // bad_request (400)
	1657: nvterr.ErrBadRequest,          // bad_request (400)
	1658: nvterr.ErrBadRequest,          // bad_request (400)
	1659: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1660: nvterr.ErrValueTooLong,        // value_too_long (422)
	1680: nvterr.ErrValueTooLong,        // value_too_long (422)
	1687: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1688: nvterr.ErrValueTooLong,        // value_too_long (422)
	1690: nvterr.ErrValueOutOfRange,     // value_out_of_range (422)
	1691: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1697: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1701: nvterr.ErrConflict,            // conflict (409)
	1706: nvterr.ErrConflict,            // conflict (409)
	1709: nvterr.ErrValueTooLong,        // value_too_long (422)
	1726: nvterr.ErrNotSupported,        // not_supported (422)
	1731: nvterr.ErrBadRequest,          // bad_request (400)
	1732: nvterr.ErrNotSupported,        // not_supported (422)
	1733: nvterr.ErrNotSupported,        // not_supported (422)
	1734: nvterr.ErrBadRequest,          // bad_request (400)
	1735: nvterr.ErrNotFound,            // not_found (404)
	1736: nvterr.ErrBadRequest,          // bad_request (400)
	1737: nvterr.ErrPreconditionFailed,  // precondition_failed (412)
	1740: nvterr.ErrNotSupported,        // not_supported (422)
	1741: nvterr.ErrNotFound,            // not_found (404)
	1742: nvterr.ErrValueTooLong,        // value_too_long (422)
	1746: nvterr.ErrBadRequest,          // bad_request (400)
	1747: nvterr.ErrBadRequest,          // bad_request (400)
	1748: nvterr.ErrValueOutOfRange,     // value_out_of_range (422)
	1757: nvterr.ErrNotSupported,        // not_supported (422)
	1758: nvterr.ErrValueOutOfRange,     // value_out_of_range (422)
	1761: nvterr.ErrDuplicateValue,      // duplicate_value (422)
	1762: nvterr.ErrDuplicateValue,      // duplicate_value (422)
	1764: nvterr.ErrBadRequest,          // bad_request (400)
	1766: nvterr.ErrInvalidState,        // invalid_state (409)
	1769: nvterr.ErrNotSupported,        // not_supported (422)
	1791: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1792: nvterr.ErrInvalidState,        // invalid_state (409)
	1793: nvterr.ErrValueTooLong,        // value_too_long (422)
	1795: nvterr.ErrNotSupported,        // not_supported (422)
	1796: nvterr.ErrNotSupported,        // not_supported (422)
	1797: nvterr.ErrBadRequest,          // bad_request (400)
	1798: nvterr.ErrBadRequest,          // bad_request (400)
	1800: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1801: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1818: nvterr.ErrInvalidLength,       // invalid_length (422)
	1819: nvterr.ErrValidationFailed,    // validation_failed (422)
	1821: nvterr.ErrPreconditionFailed,  // precondition_failed (412)
	1822: nvterr.ErrPreconditionFailed,  // precondition_failed (412)
	1824: nvterr.ErrPreconditionFailed,  // precondition_failed (412)
	1825: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1826: nvterr.ErrAlreadyExists,       // already_exists (409)
	1827: nvterr.ErrInvalidFormat,       // invalid_format (422)
	1828: nvterr.ErrConflict,            // conflict (409)
	1829: nvterr.ErrConflict,            // conflict (409)
	1830: nvterr.ErrConflict,            // conflict (409)
	1831: nvterr.ErrAlreadyExists,       // already_exists (409)
	1832: nvterr.ErrConflict,            // conflict (409)
	1833: nvterr.ErrConflict,            // conflict (409)
	1834: nvterr.ErrConflict,            // conflict (409)
	1845: nvterr.ErrNotSupported,        // not_supported (422)
	1846: nvterr.ErrNotSupported,        // not_supported (422)
	1847: nvterr.ErrNotSupported,        // not_supported (422)
	1848: nvterr.ErrNotSupported,        // not_supported (422)
	1849: nvterr.ErrNotSupported,        // not_supported (422)
	1850: nvterr.ErrNotSupported,        // not_supported (422)
	1851: nvterr.ErrNotSupported,        // not_supported (422)
	1852: nvterr.ErrNotSupported,        // not_supported (422)
	1853: nvterr.ErrNotSupported,        // not_supported (422)
	1854: nvterr.ErrNotSupported,        // not_supported (422)
	1855: nvterr.ErrNotSupported,        // not_supported (422)
	1856: nvterr.ErrNotSupported,        // not_supported (422)
	1857: nvterr.ErrNotSupported,        // not_supported (422)
	1859: nvterr.ErrDuplicateValue,      // duplicate_value (422)
	1860: nvterr.ErrValueTooLong,        // value_too_long (422)
	1861: nvterr.ErrNotSupported,        // not_supported (422)
	1865: nvterr.ErrNotSupported,        // not_supported (422)
	1869: nvterr.ErrConflict,            // conflict (409)
	1901: nvterr.ErrNotSupported,        // not_supported (422)
	1903: nvterr.ErrNotSupported,        // not_supported (422)
	1904: nvterr.ErrNotSupported,        // not_supported (422)
	1905: nvterr.ErrNotSupported,        // not_supported (422)
	1906: nvterr.ErrValueNotAllowed,     // value_not_allowed (422)
	1907: nvterr.ErrNotSupported,        // not_supported (422)
	1910: nvterr.ErrNotSupported,        // not_supported (422)
	1911: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1912: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1916: nvterr.ErrValueOutOfRange,     // value_out_of_range (422)
	1917: nvterr.ErrValueOutOfRange,     // value_out_of_range (422)
	1918: nvterr.ErrInvalidDate,         // invalid_date (422)
	1919: nvterr.ErrInvalidFormat,       // invalid_format (422)
	1921: nvterr.ErrInvalidFormat,       // invalid_format (422)
	1922: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1927: nvterr.ErrDatabaseUnavailable, // database_unavailable (503)
	1932: nvterr.ErrNotFound,            // not_found (404)
	1933: nvterr.ErrBadRequest,          // bad_request (400)
	1934: nvterr.ErrAlreadyExists,       // already_exists (409)
	1957: nvterr.ErrNotFound,            // not_found (404)
	1958: nvterr.ErrInvalidFormat,       // invalid_format (422)
	1959: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	1965: nvterr.ErrBadRequest,          // bad_request (400)
	1967: nvterr.ErrValueTooLarge,       // value_too_large (422)
	1968: nvterr.ErrAlreadyExists,       // already_exists (409)
	1969: nvterr.ErrDatabaseTimeout,     // database_timeout (504)
	1970: nvterr.ErrNotSupported,        // not_supported (422)
	1971: nvterr.ErrNotSupported,        // not_supported (422)
	1973: nvterr.ErrAlreadyExists,       // already_exists (409)
	1974: nvterr.ErrNotFound,            // not_found (404)
	1975: nvterr.ErrAlreadyExists,       // already_exists (409)
	1976: nvterr.ErrNotFound,            // not_found (404)
	1977: nvterr.ErrInvalidFormat,       // invalid_format (422)
	1978: nvterr.ErrInvalidParameter,    // invalid_parameter (400)
	3007: nvterr.ErrBadRequest,          // ER_CANNOT_DISCARD_TEMPORARY_TABLE
	3011: nvterr.ErrNotFound,            // ER_REFERENCED_TRG_DOES_NOT_EXIST_MYSQL
	3012: nvterr.ErrNotSupported,        // ER_EXPLAIN_NOT_SUPPORTED
	3013: nvterr.ErrInvalidParameter,    // ER_INVALID_FIELD_SIZE
	3014: nvterr.ErrNotSupported,        // ER_MISSING_HA_CREATE_OPTION
	3016: nvterr.ErrBadRequest,          // ER_PASSWORD_EXPIRE_ANONYMOUS_USER
	3018: nvterr.ErrNotSupported,        // ER_NO_FT_MATERIALIZED_SUBQUERY
	3020: nvterr.ErrInvalidParameter,    // ER_INVALID_ARGUMENT_FOR_LOGARITHM
	3026: nvterr.ErrDuplicateValue,      // ER_DUP_LIST_ENTRY
	3027: nvterr.ErrBadRequest,          // ER_SQL_MODE_NO_EFFECT
	3028: nvterr.ErrBadRequest,          // ER_AGGREGATE_ORDER_FOR_UNION
	3029: nvterr.ErrBadRequest,          // ER_AGGREGATE_ORDER_NON_AGG_QUERY
	3032: nvterr.ErrDatabaseUnavailable, // ER_SERVER_OFFLINE_MODE
	3033: nvterr.ErrInvalidParameter,    // ER_GIS_DIFFERENT_SRIDS
	3034: nvterr.ErrNotSupported,        // ER_GIS_UNSUPPORTED_ARGUMENT
	3037: nvterr.ErrInvalidParameter,    // ER_GIS_INVALID_DATA
	3038: nvterr.ErrRequired,            // ER_BOOST_GEOMETRY_EMPTY_INPUT_EXCEPTION
	3039: nvterr.ErrInvalidParameter,    // ER_BOOST_GEOMETRY_CENTROID_EXCEPTION
	3040: nvterr.ErrInvalidParameter,    // ER_BOOST_GEOMETRY_OVERLAY_INVALID_INPUT_EXCEPTION
	3042: nvterr.ErrInvalidParameter,    // ER_BOOST_GEOMETRY_SELF_INTERSECTION_POINT_EXCEPTION
	3045: nvterr.ErrValueOutOfRange,     // ER_STD_DOMAIN_ERROR
	3046: nvterr.ErrInvalidLength,       // ER_STD_LENGTH_ERROR
	3047: nvterr.ErrInvalidParameter,    // ER_STD_INVALID_ARGUMENT
	3048: nvterr.ErrValueOutOfRange,     // ER_STD_OUT_OF_RANGE_ERROR
	3049: nvterr.ErrValueOutOfRange,     // ER_STD_OVERFLOW_ERROR
	3050: nvterr.ErrValueOutOfRange,     // ER_STD_RANGE_ERROR
	3051: nvterr.ErrValueOutOfRange,     // ER_STD_UNDERFLOW_ERROR
	3055: nvterr.ErrInvalidFormat,       // ER_GIS_DATA_WRONG_ENDIANESS
	3057: nvterr.ErrInvalidParameter,    // ER_USER_LOCK_WRONG_NAME
	3060: nvterr.ErrNotSupported,        // ER_ALTER_OPERATION_NOT_SUPPORTED_REASON_GIS
	4002: nvterr.ErrBadRequest,          // ER_WITH_COL_WRONG_LIST
	4003: nvterr.ErrBadRequest,          // ER_TOO_MANY_DEFINITIONS_IN_WITH_CLAUSE
	4004: nvterr.ErrDuplicateValue,      // ER_DUP_QUERY_NAME
	4005: nvterr.ErrBadRequest,          // ER_RECURSIVE_WITHOUT_ANCHORS
	4006: nvterr.ErrBadRequest,          // ER_UNACCEPTABLE_MUTUAL_RECURSION
	4007: nvterr.ErrBadRequest,          // ER_REF_TO_RECURSIVE_WITH_TABLE_IN_DERIVED
	4008: nvterr.ErrBadRequest,          // ER_NOT_STANDARD_COMPLIANT_RECURSIVE
	4009: nvterr.ErrNotFound,            // ER_WRONG_WINDOW_SPEC_NAME
	4010: nvterr.ErrDuplicateValue,      // ER_DUP_WINDOW_NAME
	4011: nvterr.ErrBadRequest,          // ER_PARTITION_LIST_IN_REFERENCING_WINDOW_SPEC
	4012: nvterr.ErrBadRequest,          // ER_ORDER_LIST_IN_REFERENCING_WINDOW_SPEC
	4013: nvterr.ErrBadRequest,          // ER_WINDOW_FRAME_IN_REFERENCED_WINDOW_SPEC
	4014: nvterr.ErrBadRequest,          // ER_BAD_COMBINATION_OF_WINDOW_FRAME_BOUND_SPECS
	4015: nvterr.ErrBadRequest,          // ER_WRONG_PLACEMENT_OF_WINDOW_FUNCTION
	4016: nvterr.ErrBadRequest,          // ER_WINDOW_FUNCTION_IN_WINDOW_SPEC
	4017: nvterr.ErrBadRequest,          // ER_NOT_ALLOWED_WINDOW_FRAME
	4018: nvterr.ErrBadRequest,          // ER_NO_ORDER_LIST_IN_WINDOW_SPEC
	4019: nvterr.ErrBadRequest,          // ER_RANGE_FRAME_NEEDS_SIMPLE_ORDERBY
	4020: nvterr.ErrInvalidParameter,    // ER_WRONG_TYPE_FOR_ROWS_FRAME
	4021: nvterr.ErrInvalidParameter,    // ER_WRONG_TYPE_FOR_RANGE_FRAME
	4022: nvterr.ErrNotSupported,        // ER_FRAME_EXCLUSION_NOT_SUPPORTED
	4023: nvterr.ErrBadRequest,          // ER_WINDOW_FUNCTION_DONT_HAVE_FRAME
	4024: nvterr.ErrInvalidParameter,    // ER_INVALID_NTILE_ARGUMENT
	4025: nvterr.ErrValidationFailed,    // ER_CONSTRAINT_FAILED
	4026: nvterr.ErrValueTooLarge,       // ER_EXPRESSION_IS_TOO_BIG
	4027: nvterr.ErrBadRequest,          // ER_ERROR_EVALUATING_EXPRESSION
	4028: nvterr.ErrBadRequest,          // ER_CALCULATING_DEFAULT_VALUE
	4029: nvterr.ErrBadRequest,          // ER_EXPRESSION_REFERS_TO_UNINIT_FIELD
	4030: nvterr.ErrBadRequest,          // ER_PARTITION_DEFAULT_ERROR
	4031: nvterr.ErrNotFound,            // ER_REFERENCED_TRG_DOES_NOT_EXIST
	4032: nvterr.ErrInvalidParameter,    // ER_INVALID_DEFAULT_PARAM
	4035: nvterr.ErrInvalidFormat,       // ER_JSON_BAD_CHR
	4036: nvterr.ErrInvalidFormat,       // ER_JSON_NOT_JSON_CHR
	4037: nvterr.ErrInvalidFormat,       // ER_JSON_EOS
	4038: nvterr.ErrInvalidFormat,       // ER_JSON_SYNTAX
	4039: nvterr.ErrInvalidFormat,       // ER_JSON_ESCAPING
	4040: nvterr.ErrValueOutOfRange,     // ER_JSON_DEPTH
	4041: nvterr.ErrInvalidFormat,       // ER_JSON_PATH_EOS
	4042: nvterr.ErrInvalidFormat,       // ER_JSON_PATH_SYNTAX
	4043: nvterr.ErrValueOutOfRange,     // ER_JSON_PATH_DEPTH
	4044: nvterr.ErrInvalidParameter,    // ER_JSON_PATH_NO_WILDCARD
	4045: nvterr.ErrInvalidParameter,    // ER_JSON_PATH_ARRAY
	4046: nvterr.ErrInvalidParameter,    // ER_JSON_ONE_OR_ALL
	4047: nvterr.ErrNotSupported,        // ER_UNSUPPORTED_COMPRESSED_TABLE
	4048: nvterr.ErrInvalidFormat,       // ER_GEOJSON_INCORRECT
	4049: nvterr.ErrInvalidFormat,       // ER_GEOJSON_TOO_FEW_POINTS
	4050: nvterr.ErrInvalidFormat,       // ER_GEOJSON_NOT_CLOSED
	4051: nvterr.ErrInvalidParameter,    // ER_JSON_PATH_EMPTY
	4053: nvterr.ErrNotSupported,        // ER_FLASHBACK_NOT_SUPPORTED
	4057: nvterr.ErrNotSupported,        // ER_ISOLATION_MODE_NOT_SUPPORTED
	4058: nvterr.ErrNotSupported,        // ER_ON_DUPLICATE_DISABLED
	4059: nvterr.ErrNotSupported,        // ER_UPDATES_WITH_CONSISTENT_SNAPSHOT
	4060: nvterr.ErrInvalidState,        // ER_ROLLBACK_ONLY
	4061: nvterr.ErrInvalidState,        // ER_ROLLBACK_TO_SAVEPOINT
	4062: nvterr.ErrNotSupported,        // ER_ISOLATION_LEVEL_WITH_CONSISTENT_SNAPSHOT
	4063: nvterr.ErrNotSupported,        // ER_UNSUPPORTED_COLLATION
	4066: nvterr.ErrInvalidFormat,       // ER_RDB_TTL_DURATION_FORMAT
	4069: nvterr.ErrNotSupported,        // ER_RDB_TTL_UNSUPPORTED
	4070: nvterr.ErrInvalidFormat,       // ER_RDB_TTL_COL_FORMAT
	4071: nvterr.ErrNotSupported,        // ER_PER_INDEX_CF_DEPRECATED
	4074: nvterr.ErrBadRequest,          // ER_SUM_FUNC_WITH_WINDOW_FUNC_AS_ARG
	4076: nvterr.ErrInvalidFormat,       // ER_GEOJSON_EMPTY_COORDINATES
	4077: nvterr.ErrNotSupported,        // ER_MYROCKS_CANT_NOPAD_COLLATION
	4078: nvterr.ErrInvalidParameter,    // ER_ILLEGAL_PARAMETER_DATA_TYPES2_FOR_OPERATION
	4079: nvterr.ErrInvalidParameter,    // ER_ILLEGAL_PARAMETER_DATA_TYPE_FOR_OPERATION
	4080: nvterr.ErrBadRequest,          // ER_WRONG_PARAMCOUNT_TO_CURSOR
	4081: nvterr.ErrNotFound,            // ER_UNKNOWN_STRUCTURED_VARIABLE
	4082: nvterr.ErrNotFound,            // ER_ROW_VARIABLE_DOES_NOT_HAVE_FIELD
	4083: nvterr.ErrBadRequest,          // ER_END_IDENTIFIER_DOES_NOT_MATCH
	4084: nvterr.ErrValueOutOfRange,     // ER_SEQUENCE_RUN_OUT
	4085: nvterr.ErrInvalidParameter,    // ER_SEQUENCE_INVALID_DATA
	4086: nvterr.ErrBadRequest,          // ER_SEQUENCE_INVALID_TABLE_STRUCTURE
	4087: nvterr.ErrBadRequest,          // ER_SEQUENCE_ACCESS_ERROR
	4089: nvterr.ErrBadRequest,          // ER_NOT_SEQUENCE
	4090: nvterr.ErrBadRequest,          // ER_NOT_SEQUENCE2
	4091: nvterr.ErrNotFound,            // ER_UNKNOWN_SEQUENCES
	4092: nvterr.ErrNotFound,            // ER_UNKNOWN_VIEW
	4093: nvterr.ErrBadRequest,          // ER_WRONG_INSERT_INTO_SEQUENCE
	4095: nvterr.ErrBadRequest,          // ER_PACKAGE_ROUTINE_IN_SPEC_NOT_DEFINED_IN_BODY
	4096: nvterr.ErrBadRequest,          // ER_PACKAGE_ROUTINE_FORWARD_DECLARATION_NOT_DEFINED
	4097: nvterr.ErrNotSupported,        // ER_COMPRESSED_COLUMN_USED_AS_KEY
	4098: nvterr.ErrNotFound,            // ER_UNKNOWN_COMPRESSION_METHOD
	4099: nvterr.ErrBadRequest,          // ER_WRONG_NUMBER_OF_VALUES_IN_TVC
	4100: nvterr.ErrBadRequest,          // ER_FIELD_REFERENCE_IN_TVC
	4101: nvterr.ErrInvalidParameter,    // ER_WRONG_TYPE_FOR_PERCENTILE_FUNC
	4102: nvterr.ErrInvalidParameter,    // ER_ARGUMENT_NOT_CONSTANT
	4103: nvterr.ErrValueOutOfRange,     // ER_ARGUMENT_OUT_OF_RANGE
	4104: nvterr.ErrInvalidParameter,    // ER_WRONG_TYPE_OF_ARGUMENT
	4105: nvterr.ErrBadRequest,          // ER_NOT_AGGREGATE_FUNCTION
	4106: nvterr.ErrBadRequest,          // ER_INVALID_AGGREGATE_FUNCTION
	4107: nvterr.ErrInvalidParameter,    // ER_INVALID_VALUE_TO_LIMIT
	4108: nvterr.ErrBadRequest,          // ER_INVISIBLE_NOT_NULL_WITHOUT_DEFAULT
	4109: nvterr.ErrNotSupported,        // ER_UPDATE_INFO_WITH_SYSTEM_VERSIONING
	4110: nvterr.ErrInvalidParameter,    // ER_VERS_FIELD_WRONG_TYPE
	4111: nvterr.ErrNotSupported,        // ER_VERS_ENGINE_UNSUPPORTED
	4112: nvterr.ErrInvalidParameter,    // ER_WRONG_VERSIONING_RANGE
	4113: nvterr.ErrInvalidParameter,    // ER_PARTITION_WRONG_TYPE
	4116: nvterr.ErrNotSupported,        // ER_VERS_DROP_PARTITION_INTERVAL
	4119: nvterr.ErrNotSupported,        // ER_VERS_ALTER_NOT_ALLOWED
	4120: nvterr.ErrNotSupported,        // ER_VERS_ALTER_ENGINE_PROHIBITED
	4121: nvterr.ErrNotSupported,        // ER_VERS_RANGE_PROHIBITED
	4122: nvterr.ErrBadRequest,          // ER_CONFLICTING_FOR_SYSTEM_TIME
	4123: nvterr.ErrBadRequest,          // ER_VERS_TABLE_MUST_HAVE_COLUMNS
	4124: nvterr.ErrBadRequest,          // ER_VERS_NOT_VERSIONED
	4125: nvterr.ErrBadRequest,          // ER_MISSING
	4126: nvterr.ErrBadRequest,          // ER_VERS_PERIOD_COLUMNS
	4127: nvterr.ErrInvalidParameter,    // ER_PART_WRONG_VALUE
	4128: nvterr.ErrBadRequest,          // ER_VERS_WRONG_PARTS
	4129: nvterr.ErrNotSupported,        // ER_VERS_NO_TRX_ID
	4130: nvterr.ErrBadRequest,          // ER_VERS_ALTER_SYSTEM_FIELD
	4131: nvterr.ErrNotSupported,        // ER_DROP_VERSIONING_SYSTEM_TIME_PARTITION
	4132: nvterr.ErrNotSupported,        // ER_VERS_DB_NOT_SUPPORTED
	4134: nvterr.ErrDuplicateValue,      // ER_VERS_DUPLICATE_ROW_START_END
	4135: nvterr.ErrAlreadyExists,       // ER_VERS_ALREADY_VERSIONED
	4137: nvterr.ErrNotSupported,        // ER_VERS_NOT_SUPPORTED
	4138: nvterr.ErrNotSupported,        // ER_VERS_TRX_PART_HISTORIC_ROW_NOT_SUPPORTED
	4140: nvterr.ErrBadRequest,          // ER_UPDATED_COLUMN_ONLY_ONCE
	4141: nvterr.ErrBadRequest,          // ER_EMPTY_ROW_IN_TVC
	4142: nvterr.ErrNotSupported,        // ER_VERS_QUERY_IN_PARTITION
	4143: nvterr.ErrNotSupported,        // ER_KEY_DOESNT_SUPPORT
	4144: nvterr.ErrNotSupported,        // ER_ALTER_OPERATION_TABLE_OPTIONS_NEED_REBUILD
	4145: nvterr.ErrLocked,              // ER_BACKUP_LOCK_IS_ACTIVE
	4146: nvterr.ErrInvalidState,        // ER_BACKUP_NOT_RUNNING
	4147: nvterr.ErrInvalidState,        // ER_BACKUP_WRONG_STAGE
	4149: nvterr.ErrInvalidParameter,    // ER_BACKUP_UNKNOWN_STAGE
	4150: nvterr.ErrLocked,              // ER_USER_IS_BLOCKED
	4151: nvterr.ErrLocked,              // ER_ACCOUNT_HAS_BEEN_LOCKED
	4152: nvterr.ErrNotSupported,        // ER_PERIOD_TEMPORARY_NOT_ALLOWED
	4153: nvterr.ErrBadRequest,          // ER_PERIOD_TYPES_MISMATCH
	4154: nvterr.ErrBadRequest,          // ER_MORE_THAN_ONE_PERIOD
	4155: nvterr.ErrInvalidParameter,    // ER_PERIOD_FIELD_WRONG_ATTRIBUTES
	4156: nvterr.ErrNotFound,            // ER_PERIOD_NOT_FOUND
	4157: nvterr.ErrBadRequest,          // ER_PERIOD_COLUMNS_UPDATED
	4158: nvterr.ErrBadRequest,          // ER_PERIOD_CONSTRAINT_DROP
	4159: nvterr.ErrValueTooLong,        // ER_TOO_LONG_KEYPART
	4160: nvterr.ErrValueTooLong,        // ER_TOO_LONG_DATABASE_COMMENT
	4161: nvterr.ErrInvalidParameter,    // ER_UNKNOWN_DATA_TYPE
	4162: nvterr.ErrBadRequest,          // ER_UNKNOWN_OPERATOR
	4163: nvterr.ErrDuplicateValue,      // ER_SP_DUP_DECL
	4164: nvterr.ErrValueOutOfRange,     // ER_PART_STARTS_BEYOND_INTERVAL
	4169: nvterr.ErrBadRequest,          // ER_NO_AUTOINCREMENT_WITH_UNIQUE
	4170: nvterr.ErrBadRequest,          // ER_KEY_CONTAINS_PERIOD_FIELDS
	4171: nvterr.ErrBadRequest,          // ER_KEY_CANT_HAVE_WITHOUT_OVERLAPS
	4172: nvterr.ErrBadRequest,          // ER_NOT_ALLOWED_IN_THIS_CONTEXT
	4174: nvterr.ErrBadRequest,          // ER_PK_INDEX_CANT_BE_IGNORED
	4176: nvterr.ErrInvalidParameter,    // ER_JSON_TABLE_ERROR_ON_FIELD
	4177: nvterr.ErrBadRequest,          // ER_JSON_TABLE_ALIAS_REQUIRED
	4178: nvterr.ErrInvalidParameter,    // ER_JSON_TABLE_SCALAR_EXPECTED
	4179: nvterr.ErrBadRequest,          // ER_JSON_TABLE_MULTIPLE_MATCHES
	4180: nvterr.ErrBadRequest,          // ER_WITH_TIES_NEEDS_ORDER
	4184: nvterr.ErrNotSupported,        // ER_PARTITION_CONVERT_SUBPARTITIONED
	4187: nvterr.ErrBadRequest,          // ER_SF_OUT_INOUT_ARG_NOT_ALLOWED
	4192: nvterr.ErrInvalidParameter,    // ER_SLAVE_STATEMENT_TIMEOUT
	4193: nvterr.ErrNotSupported,        // ER_JSON_INVALID_VALUE_FOR_KEYWORD
	4194: nvterr.ErrBadRequest,          // ER_JSON_SCHEMA_KEYWORD_UNSUPPORTED
	4195: nvterr.ErrBadRequest,          // ER_JSON_NO_VARIABLE_SCHEMA
	4196: nvterr.ErrBadRequest,          // ER_SEQUENCE_TABLE_HAS_WRONG_NUMBER_OF_COLUMNS
	4197: nvterr.ErrBadRequest,          // ER_SEQUENCE_TABLE_CANNOT_HAVE_ANY_KEYS
	4198: nvterr.ErrBadRequest,          // ER_SEQUENCE_TABLE_CANNOT_HAVE_ANY_CONSTRAINTS
	4199: nvterr.ErrBadRequest,          // ER_SEQUENCE_TABLE_ORDER_BY
	4201: nvterr.ErrBadRequest,          // ER_INCORRECT_COLUMN_NAME_COUNT
	4203: nvterr.ErrInvalidFormat,       // ER_VECTOR_BINARY_FORMAT_INVALID
	4204: nvterr.ErrInvalidFormat,       // ER_VECTOR_FORMAT_INVALID
	4206: nvterr.ErrInvalidParameter,    // ER_VEC_DISTANCE_TYPE
	4212: nvterr.ErrNotImplemented,      // ER_NOT_IMPLEMENTED_FOR_GEOGRAPHIC_SRS
	4213: nvterr.ErrValueOutOfRange,     // ER_LONGITUDE_OUT_OF_RANGE
	4214: nvterr.ErrValueOutOfRange,     // ER_LATITUDE_OUT_OF_RANGE
	4215: nvterr.ErrBadRequest,          // ER_GIS_DIFFERENT_SRIDS_AGGREGATION
	4218: nvterr.ErrBadRequest,          // ER_WARN_OPTIMIZER_HINT_SYNTAX_ERROR
	4219: nvterr.ErrBadRequest,          // ER_WARN_CONFLICTING_HINT
	4220: nvterr.ErrBadRequest,          // ER_WARN_UNKNOWN_QB_NAME
	4221: nvterr.ErrBadRequest,          // ER_UNRESOLVED_TABLE_HINT_NAME
	4222: nvterr.ErrBadRequest,          // ER_UNRESOLVED_INDEX_HINT_NAME
	4223: nvterr.ErrNotSupported,        // ER_HINTS_INSIDE_VIEWS_NOT_SUPPORTED
	4224: nvterr.ErrBadRequest,          // ER_WARN_MALFORMED_HINT
	4225: nvterr.ErrBadRequest,          // ER_WARN_HINTS_ON_INSERT_PART_OF_INSERT_SELECT
	4227: nvterr.ErrConflict,            // ER_ACTIVE_CONNECTIONS_FOR_USER_TO_DROP
	4228: nvterr.ErrNotFound,            // ER_ASSOC_ARRAY_ELEM_NOT_FOUND
	4229: nvterr.ErrInvalidParameter,    // ER_NULL_FOR_ASSOC_ARRAY_INDEX
	4230: nvterr.ErrBadRequest,          // ER_NEED_NAMED_ASSOCIATION
	4231: nvterr.ErrBadRequest,          // ER_INVALID_USE_OF_ORA_JOIN
	4232: nvterr.ErrBadRequest,          // ER_INVALID_USE_OF_ORA_JOIN_CYCLE
	4233: nvterr.ErrBadRequest,          // ER_INVALID_USE_OF_ORA_JOIN_OUTER_REF
	4234: nvterr.ErrBadRequest,          // ER_INVALID_USE_OF_ORA_JOIN_ONE_TABLE
	4235: nvterr.ErrBadRequest,          // ER_INVALID_USE_OF_ORA_JOIN_WRONG_FUNC
	4236: nvterr.ErrBadRequest,          // ER_INVALID_USE_OF_ORA_JOIN_MIX
	4238: nvterr.ErrBadRequest,          // ER_WARN_CONFLICTING_INDEX_HINT_FOR_TABLE
	4239: nvterr.ErrBadRequest,          // ER_WARN_CONFLICTING_INDEX_HINT_FOR_KEY
	4240: nvterr.ErrBadRequest,          // ER_WARN_CONFLICTING_COMPOUND_INDEX_HINT_FOR_TABLE
	4241: nvterr.ErrBadRequest,          // ER_WARN_CONFLICTING_COMPOUND_INDEX_HINT_FOR_KEY
	4242: nvterr.ErrBadRequest,          // ER_WARN_NO_IMPLICIT_QB_NAMES_IN_VIEW
	4255: nvterr.ErrNotSupported,        // ER_UNSUPPORTED_DATA_TYPE_ATTRIBUTE
	4257: nvterr.ErrBadRequest,          // ER_CANNOT_CAST_ON_IDENT1_ASSIGNMENT_FOR_OPERATION
	4258: nvterr.ErrBadRequest,          // ER_CANNOT_CAST_ON_IDENT2_ASSIGNMENT_FOR_OPERATION
	4259: nvterr.ErrBadRequest,          // ER_WARN_AMBIGUOUS_QB_NAME
	4260: nvterr.ErrBadRequest,          // ER_WARN_IMPLICIT_QB_NAME_FOR_UNION
	4261: nvterr.ErrBadRequest,          // ER_WARN_QB_NAME_PATH_INVALID_SELECT
	4262: nvterr.ErrBadRequest,          // ER_WARN_QB_NAME_PATH_SELECT_NOT_FOUND
	4263: nvterr.ErrBadRequest,          // ER_WARN_QB_NAME_PATH_VIEW_NOT_FOUND
	4264: nvterr.ErrNotSupported,        // ER_WARN_QB_NAME_PATH_NOT_SUPPORTED_INSIDE_VIEW
	4265: nvterr.ErrBadRequest,          // ER_PARTITION_INTERVAL_NOT_LIST
	4266: nvterr.ErrBadRequest,          // ER_PARTITION_INTERVAL_FINER_THAN_DATE
	4267: nvterr.ErrBadRequest,          // ER_PARTITION_INTERVAL_MAXVALUE
	4268: nvterr.ErrBadRequest,          // ER_RANGE_INTERVAL_PART_FAILED
	4276: nvterr.ErrBadRequest,          // ER_DDL_ILLEGAL_FK
	4279: nvterr.ErrInvalidParameter,    // ER_INVALID_HOST_NETMASK
	4282: nvterr.ErrBadRequest,          // ER_CANT_SET_IN_SUBSTATEMENT
	5002: nvterr.ErrNotSupported,        // CR_ASYNC_NOT_SUPPORTED
	5003: nvterr.ErrNotSupported,        // CR_FUNCTION_NOT_SUPPORTED
	5004: nvterr.ErrNotFound,            // CR_FILE_NOT_FOUND
	5006: nvterr.ErrInvalidParameter,    // CR_BULK_WITHOUT_PARAMETERS
	5007: nvterr.ErrBadRequest,          // CR_INVALID_STMT
	5009: nvterr.ErrInvalidParameter,    // CR_INVALID_PARAMETER
	5024: nvterr.ErrInvalidParameter,    // CR_INVALID_CLIENT_FLAG
	5025: nvterr.ErrInvalidState,        // CR_STMT_NO_RESULT
	5027: nvterr.ErrValueOutOfRange,     // CR_ERR_TOO_MUCH_COLUMNS
}
