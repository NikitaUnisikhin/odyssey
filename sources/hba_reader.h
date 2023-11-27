#ifndef ODYSSEY_HBA_READER_H
#define ODYSSEY_HBA_READER_H

int od_hba_reader_parse(od_config_reader_t *reader);
int od_hba_reader_prefix(od_hba_rule_t *hba, char *prefix);
int od_hba_reader_address(struct sockaddr_storage *dest,
			  const char *addr);
void od_hba_reader_error(od_config_reader_t *reader, char *msg);

#endif /* ODYSSEY_HBA_READER_H */
