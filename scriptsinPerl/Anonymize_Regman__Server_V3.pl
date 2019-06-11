#!/usr/bin/perl 
# Active State Perl   v5.12  for i386 su piattaforma Linux
# Data:         07/01/2019
# Versione      1.0 Anonimizzazione dati per PoC Machine Learning con Google/Ingenia - Prima Versione - Dati AVS
# Esempio riga da file sorgente
# AVS
#       0                       1                2       3              4                               5                       6               7                 8   9     10                      11                              12                      13              14                                15    16                 17
# Dispositivo Timestamp CLI CPE-ID Tipo Contenuto ID contenuto ??? Tipo Acquisto null null Account Account Conciliato Piattaforma Vendor Altro tipo contenuto null Tipo linea TEMPO DI VISIONE (secondi)
#             0                 1                        2                  3           4           5         6        7         8   9     10                11       12           13         14                     15   16          17                         
# CUBOTECHNICOLOR_ANDROIDTV     2018-11-30T00:00:08 000481808957 A7655449000686B9 CATALOGUE  50707453 STANDARD ABBONAMENTO |   |   |g.korenjak@alice.it|   |  CUBOVISION  |  ANDROID | Abbonamento;DownloadPlay |  | RETE FISSA | 48
# Esempio riga CDN
#           0                                                   1                                  2                                    3                               5                                   6                                                                                                                          8
# time-recvd(millisecond)       time-to-serve(microsecond)      client-ip       request-desc/response-status    bytes-sent      request-method  request-url        mine-type        req-header(user-agent)  abr-protocol    bitrate asset-id        session-id      entry-generated-time    client-type     profile
# [30/Nov/2018:19:59:59.786+0000]       218801                          79.18.78.113    TCP_MISS/200                             122508         GET                http://voddashhttps.cb.ticdn.it/videoteca2/V3/Film/2016/02/50541341/DASH/11475246/r16/191.mp4f   video/mp4       "Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko"     -       0       -       -
# Esempio riga NGASP
# cpeid; tgu;trap_timestamp; deviceid; devicetype; mode; originipaddress; averagebitrate; avgsskbps; bufferingduration; callerclass; callerrorcode; callerrormessage; callerrortype; callurl; errordesc; errorreason; eventname; levelbitrates; linespeedkbps; maxsschunkkbps; maxsskbps; minsskbps; streamingtype; videoduration; videoposition; videotitle; videotype; videourl; eventtype; fwversion; networktype; ra_version; update_time; trap_provider; mid; service_id; service_id_version; date_rif;   video_provider; max_upstream_net_latency; min_upstream_net_latency; avg_upstream_net_latency; max_downstream_net_latency; min_downstream_net_latency; avg_downstream_net_latency; max_platform_latency; min_platform_latency; avg_platform_latency; packet_loss; preloaded_app_v 
# Autore:       Giovanni Lepri
# Nome:                 Anonymize.pl
# Scopo:        Anonimizzazione dati per trial machine learning google
#
#++++++++++++++++++++++++++++
# Librerie
#++++++++++++++++++++++++++++ 
use IO::Uncompress::Unzip;                                                                                                                      # Libreria di deccompressione zip
use IO::Compress::Zip;                                                                                                                          # Libreria di compressione zip
#
use Compress::Zlib;                                                                                                                                     # Libreria di compressione gz
#
use Crypt::Lite;                                                                                                                                        # Libreria di criptazione di stringhe
# use Switch;                                                                                                                                           # Implementa lo statement Case - Switch
#++++++++++++++++++++++++++++
use constant STEP => 10000;                                                                                                                     # Passo di scrittura su riga di avanzamento
#++++++++++++++++++++++++++++
# Definizione dei campi di da trascrivere sul file criptato in funzione della tipologia. I valori da riportare in uscita sono quelli maggiori di zero.
# Attualmente non utilizzato
#
%RecordType = {
        'CONTENT_CRIPTED_ERROR', '2',
        'CONTENT_PLAYBACK_ERROR', '2',
        'END_BUFFERING', '2',
        'NETWORK_ERROR', '2',
        'SERVER_ERROR', '2',
        'SS_QUALITY', '2',
        'T2FP', '2',
        'TtFP', '2',
        'APPLICATION_BOOT', '0',
        'BUFFER_UNDERFLOW', '1',
        'CAROUSEL_OK', '1',
        'DEVICE_BOOT', '0',
        'DTT_STATUS', '0',
        'END_BUFFERING', '1',
        'ETHERNET', '0',
        'FFW', '1',
        'LONG_BUFFERING', '1',
        'PAUSE', '1',
        'PLAY', '1',
        'PLAY_REQUEST', '1',
        'PRELOADED_UPGRADE', '0',
        'REWIND', '1',
        'SPEED_TEST', '0',
        'STANDBY', '0',
        'STOP', '1',
        'SV_CLOSE', '0',
        'SV_OPEN', '0',
        'SW_UPGRADE', '1',
        'TEST_LINEA_ADSL', '0',
        'WAKEUP', '0',
        'WIFI', '0'
};
#++++++++++++++++++++++++++++
#use Date::Calc qw(Week_of_Year Day_of_Week Week_Number Day_of_Year Delta_DHMS);        # Libreria per manipolazione date
#use Carp;
#++++++++++++++++++++++++++++
# Dichiarazione subroutines


#++++++++++++++++++++++++++++
# Legge il tipo di elaborazione richiesta ed il tipo di server
my $EvalutionType = shift or die "Usage: $0 tipo di elaborazione richiesta: AVS|CDN|REGMAN host: server|Mac\n";
my $myHost = shift or die "Usage: $0 tipo di elaborazione richiesta: AVS|CDN|REGMAN host: server|Mac\n";
# Lista di parametri ammessi nella rivadi comando
%EvaluationList = ( 'AVS', '0', 'CDN', '0', 'NGASP', '1' );
# Verifica il parametro
if ($EvaluationList{$EvalutionType} == 1){
        print "\n **** Inizio anonimizzazione file $EvalutionType ****\n"
} else
{
        die " **** Tipo di elaborazione non prevista: $EvalutionType ****\n"
}
# Dichiarazione directory e nomi file
# $myWorkdir = '/Volumes/gioMacBookL/Automation/Analytics/';                                            # Directory di lavoro - MAC
$myWorkdir = '/home/gioml/Analytics/';                                                                                  # Directory di lavoro - Server
#$myAVSDataDir = '/Volumes/gioMacBookL/Automation/';                                                            # Directory dati origine  AVS e Trap - MAC
#$myTrapDataDir = '/Volumes/gioMacBookL/Automation/';                                                   # Directory dati origine  TRAP da NGASP - Mac
#$myAVSDataDir = '/home/gioml/Analytics/AVS/';                                                                  # Directory dati origine  AVS - Server
$myTrapDataDir = '/home/gioml/Analytics/REGMAN/';                                                               # Directory dati origine  TRAP da REGMAN - Server
$myCDNDataDir[0]= '/Volumes/gioMacBookL/Automation/Analytics/h_20_UTC';                 # Directory dati di origine CDN ore 20 - MAC
$myCDNDataDir[1]= '/Volumes/gioMacBookL/Automation/Analytics/h_21_UTC';                 # # Directory dati di origine CDN ore 21 -MAC
#$myCDNDataDir[0]= '/home/gioml/Analytics/h_20_UTC';                                                    # Directory dati di origine CDN ore 20 - Server
#$myCDNDataDir[1]= '/home/gioml/Analytics/h_21_UTC';                                                    # Directory dati di origine CDN ore 21 - Server
# File contenenti i dati di AVS e Trap
# $myTrapFile = $myTrapDataDir . 'testvideo.csv';                                                               # File contenente le trap dei dispositivi (STB, Connected TV)
$myTrapFile = $myTrapDataDir . 'cubo_traps_filtered_20181130.csv';                              # File contenente le trap dei dispositivi (STB, Connected TV), filtrati dei contenuti di Tim Music, rispetto all'orginale
# $myAVSFile = $myAVSDataDir . 'Fruizioni_AVS_30112018.zip';                                    # File contenente i dati di AVS (Front End di Tim Vision)
# $myAVSFile = $myAVSDataDir . 'Cartel2.csv';
# Redirige stderr e impone autoflush
my $logfile = $myWorkdir. 'logerrAnonymizeNGASP.txt';
open STDERR, ">$logfile" or die "\ **** Impossibile aprire file di errori $logfile *****\n";
local $| = 1;                                                                                                                                           # standard out autoflush
print "\n ++++++++ \nFile di log creato: $logfile ++++++++\n";
# Inizializzazione variabili per la visualizzazione dell'avanzamento attività
$av =  STEP - 1;                                        # passo di avanzamento visualizzazione 1 milione di record
$index = 0;                                                     # per record elaborati == numero di linee lette
$myCounter = 0;                                         # per testo a capo
#++++++++++++++++++++++++++++
# Costruisce nome del file di uscita e lo apre - Dati NGASP
# die ("Debug fase 1");
#
$CliStats = $myWorkdir.'NGASP_TRAPS_20181130_anonymized.txt.gz';                                                                                # i file elaborati contengono le trap dei STB e CTV ottenuti da Regman Monitoraggio
if ( -e $CliStats) {unlink $CliStats or die "\n ***** Impossibile cancellare il file $CliStats \n" };   # creerà un file di uscita con le gli IP mascherati per i chunk erogati
my $gz_out = gzopen($CliStats, "wb") or die "\n **** Impossibile aprire $CliStats in scrittura: $! ****\n";
# Scrive l'header di uscita - Dati di AVS
$PrintedChars = 0;                                                                                                                                                 # Numero di caratteri scritti nel file risultato
# Intestazione 24 campi
#                0          1   2               3         4        5       6               7               8              9         10        11            12         13         14            15           16           17          18              19         20          21         22                  23                           
$Outputline = "cpeid; tgu;trap_timestamp;devicetype;originipaddress;averagebitrate;bufferingduration;errordesc;eventname;maxsschunkkbps;maxsskbps;minsskbps;videoduration;videoposition;videotitle;streamingtype;trap_provider;fwversion;networktype;update_time;provider\n";
$PrintedChars += $gz_out->gzwrite($Outputline);
#+++++++++++++++++++++++++++
# Time stamp - Valuta le prestazioni e visualizza avanzamento
($sec,$min,$hour,$mday,$mon,$year,$wday,$yday,$isdst)=localtime(time);
$StartTime = sprintf "%4d-%02d-%02d %02d:%02d:%02d",$year+1900,$mon+1,$mday,$hour,$min,$sec;
print STDERR "\n **** $StartTime - Google PoC - Anonimizzazione dati: Fonte NGASP ****\n";
print "\n ++++++++ \n$StartTime> Avvio elaborazione Dati NGASP\n";
print "\n                               - Milioni di record valutati - \n\n";
print "1       10        20       30        40        50        60        70        80       90       100\n";   #
#+++++++++++++++++++++++++++
# Ciclo su tutti i file delle directory contenenti i dati CDN
# Variabili di contesto
# $HeaderYesNO = 0;                                                                                                                                             # Variabile sentinella per la scrittura della riga di header
$myCrypt = Crypt::Lite->new( debug => 0 );                                                                                              # Instanzia la funzione di criptazione
$mySecret = 'tim_UFE1';                                                                                                                                 # Stringa con chiave di criptazione
# Inizia il ciclo
$ListaFile{$myTrapFile} = 1;                                                                                                                            # Memorizza nell'hash il nome del file contenente le fruizioni su AVS
foreach my $FileInput (sort keys %ListaFile)
{
        # Apre il file di ingresso, i-esimo elemento della lista
    # my $gz_in = new IO::Uncompress::Unzip $FileInput [OPTS]
        open($fh, '<', $FileInput)
        or die "\n$FileInput Open failed: text file read error\n";
        # Ciclo sul file dati, riga per riga
#       if ($HeaderYesNO == 0) {$HeaderYesNO = 1; $PrintedChars += $gz_out->gzwrite($HeaderCDN);}       # Scrive l'intestazione una sola volta nel file risultato
    # LOADLOOP: while ($line = $gz_in->getline())  {
        LOADLOOP: while ($line = <$fh>) {
                # print   "$line\n";                                                                                            # debug only
                # chomp $line;
                @fields = split(';', $line);                                                                                    # suddivisione della riga in stringhe separate dal carattere pipe
                $CLI = $fields[1];                                                                                              # CLI
                # $CPE_ID = $fields[0];                                                                                                                 # CPE_ID
                # $IP_ADD = $fields[6];                                                                                                                 # IP Address
                # print "$IP_ADD, $CLI, $CPE_ID \n";                                                                                    # Debug only
                if ($CLI !~ /^[0-9]+/){ $fields[1] = 'OLO_000'; }                                                               # Verifica la presenza dell'indirizzo IP
                        # print "\n **** CLI e CPE: $fields[2] $CPE_ID\n";                                                      # debug only
                        $fields[1] = $myCrypt->encrypt($fields[1], $mySecret);                                          # Cripta il CLI
                        # print "$fields[2] ****\n";                                                                                            # debug only
                        $fields[6] = $myCrypt->encrypt($fields[6], $mySecret);                                          # Cripta l'indirizzo IP
                        $fields[0] = $myCrypt->encrypt($fields[0], $mySecret);                                          # Cripta il CPE ID
# Costruisce la riga per la scrittura nel file risultato, in funzione del tipo di evento e del dispositivo
#                0          1       2               3         4              5              6               7          8             9         10        11            12         13           14            15           16           17          18         19         20                                     
# Outputline = "cpeid; tgu;trap_timestamp;devicetype;originipaddress;averagebitrate;bufferingduration;errordesc;eventname;maxsschunkkbps;maxsskbps;minsskbps;videoduration;videoposition;videotitle;streamingtype;trap_provider;fwversion;networktype;update_time;provider\n";
#
# La descrizione dell'evento si trova in colonna 26 - Le app su TV Bravia hanno formato diverso
#
                        if (length($fields[26]) > 0 ){
                                if($fields[26] eq 'SS_QUALITY') {                                                                               # evento che contiene i dati di qualità della fruizione
                                        if ($fields[4] ne 'BRAVIA') {
                                                $Outputline = "$fields[0];$fields[1];$fields[2];$fields[4];$fields[6];$fields[9];-1;-1;$fields[26];$fields[31];$fields[32];$fields[33];-1;-1;$fields[40];$fields[41];$fields[43];$fields[44];$fields[45];$fields[48];$fields[56]\n";         # Linea criptata, esclusi i campi non di interesse
                                                # print "SS_Quality no bravia \n";
                                        } else {
                                                $Outputline = "$fields[0];$fields[1];$fields[2];$fields[4];$fields[6];$fields[9];-1;-1;$fields[26];$fields[31];$fields[32];$fields[33];-1;-1;$fields[40];-1;$fields[44];$fields[45];$fields[46];$fields[49];$fields[58]\n";  # Linea criptata, esclusi i campi non di interesse
                                        }
                                }
                                elsif ($fields[26] eq 'TtFP' || $fields[26] eq 'T2FP' ) {                               # evento che contiene il tempo di avvio della riproduzione
                                        if ($fields[4] ne 'BRAVIA') {
                                                $Outputline = "$fields[0];$fields[1];$fields[2];$fields[4];$fields[6];-1;$fields[11];-1;$fields[26];-1;-1;-1;-1;-1;$fields[43];$fields[44];$fields[46];$fields[47];$fields[48];$fields[51];$fields[59]\n";   # Linea criptata, esclusi i campi non di interesse
                                        } else {
                                                $Outputline = "$fields[0];$fields[1];$fields[2];$fields[4];$fields[6];-1;$fields[11];-1;$fields[26];-1;-1;-1;-1;-1;$fields[43];-1;$fields[47];$fields[48];$fields[49];$fields[52];$fields[61]\n";    # Linea criptata, esclusi i campi non di interesse
                                        }
                                } else {
                                        if ($fields[4] ne 'BRAVIA') {                                                                           # eventi di errore (network error, buffering, content crypt error, etc.)
                                                $Outputline = "$fields[0];$fields[1];$fields[2];$fields[4];$fields[6];-1;$fields[11];$fields[23];$fields[26];-1;-1;-1;-1;-1;$fields[43];$fields[44];$fields[46];$fields[47];$fields[48];$fields[51];$fields[59]\n";  # Linea criptata, esclusi i campi non di interesse
                                        } else {
                                                $Outputline = "$fields[0];$fields[1];$fields[2];$fields[4];$fields[6];-1;$fields[11];$fields[23];$fields[26];-1;-1;-1;-1;-1;$fields[43];-1;$fields[47];$fields[48];$fields[49];$fields[52];$fields[61]\n";   # Linea criptata, esclusi i campi non di interesse
                                        }
                                }
#
# La descrizione dell'evento si trova in colonna 27
#
                        } else {                                                                                                                                   # altri eventi
                                if ($fields[4] ne 'BRAVIA') {
                                        if ($fields[27] eq 'FFW' || $fields[27] eq 'REWIND'){                           # evento avanti veloce e ritorno indietro
                                                $Outputline = "$fields[0];$fields[1];$fields[2];$fields[4];$fields[6];-1;-1;-1;$fields[27];-1;-1;-1;-1;$fields[42];$fields[43];$fields[44];$fields[46];$fields[47];$fields[48];$fields[51];$fields[59]\n";   # Linea criptata, esclusi i campi non di interesse
                                        } elsif  ($fields[27] ne 'STOP') {                                                                      # evento di fine fruizione
                                                $Outputline = "$fields[0];$fields[1];$fields[2];$fields[4];$fields[6];-1;-1;-1;$fields[27];-1;-1;-1;-1;-1;$fields[44];$fields[45];$fields[47];$fields[48];$fields[49];$fields[52];$fields[60]\n";    # Linea criptata, esclusi i campi non di interesse
                                        } else {                                                                                                                   # evento di avvio riproduzione, pausa, etc.
                                                $Outputline = "$fields[0];$fields[1];$fields[2];$fields[4];$fields[6];-1;-1;-1;$fields[27];-1;-1;-1;$fields[40];$fields[41];$fields[42];$fields[43];$fields[45];$fields[46];$fields[47];$fields[50];$fields[58]\n";  # Linea criptata, esclusi i campi non di interesse
                                        }
                                } else {
                                                $Outputline = "$fields[0];$fields[1];$fields[2];$fields[4];$fields[6];-1;-1;-1;$fields[27];-1;-1;-1;-1;-1;$fields[44];-1;$fields[48];$fields[49];$fields[50];$fields[53];$fields[62]\n";     # Linea criptata, esclusi i campi non di interesse
                                } 
                        } 
                        # print "Linea originale \n $line\n Linea Criptata \n $Outputline \n";          # Debug only, confronto pre e post criptazione
                        $PrintedChars += $gz_out->gzwrite($Outputline)
                                or die "\n **** Impossibile scrivere sul file $CliStats ****\n";                # Scrive la riga sul file di uscita, dopo l'anonimizzazione dell'indirizzo IP
                        $ChunkCount++;                                                                                                                          # Aggiorna il progressivo del file all'interno della settimana
# Visualizza l'avanzamento
                $index++;
                if ($index > $av) {
                        $av += STEP;
                        print STDOUT "+";                                                                                                                       # debug - Togliere il commento in caso di esecuzione
                        $myCounter++;
                        if ($myCounter > 98) {
                                $myCounter = 0;
                                print "\n";
                        }
                }
        }       # Fine esplorazione di un file
        # Chiude il file di I/O
    # close($gz_in) or die "\n **** Impossibile chiudere il file $myAVSFile ****\n";
        # $gz_out->gzclose() or die "\n **** Impossibile chiudere il file $CliStats ****\n";                                                                       # Debug only
        # $gz_out->gzclose();
        # die "\n **** Esempio file cdn - $PrintedChars caratteri scritti ****\n";              # Debug only
}       # Fine ciclo sui nomi file
# +++++++++++++++++++++++++++++++
# Synthesis
print STDERR "\n **** NGASP file Evaluated. $ChunkCount record found ****\n";
print "\nTerminata elaborazione dei file NGASP - $ChunkCount record validi scritti su file risultato ****\n";
# Chiusura file di I/O --> scrittura
$gz_out->gzclose();  # or die "\n **** Impossibile chiudere il file $CliStats ****\n";
# Chiusura file di I/O --> lettura
close $fh;
# +++++++++++++++++++++++++++++++
# Time Stamp fine esecuzione
($sec,$min,$hour,$mday,$mon,$year,$wday,$yday,$isdst)=localtime(time);
$StopTime = sprintf "%4d-%02d-%02d %02d:%02d:%02d",$year+1900,$mon+1,$mday,$hour,$min,$sec;
print "\n **** Fine elaborazione alle $StopTime ****\n";
print STDERR "\n **** Fine elaborazione alle $StopTime ****\n";
# Completa l'esecuzione
exit(1);
# +++++++++++++++++++++++++++++++
sub ScanCDNDir ($$)
{
        my $localCDNDataDir = shift;                                                    # Directory contenente i file CDN, formato gzip
        my $CDNFileTemplate = shift;
        my $FullPathFileName = '';                                                              # Contiene il path e il nome del file
        my $FileCount = 0;                                                                              # Numero di file trovati
        my $UnusableFile = 0;                                                                   # Numero di file inutilizzabili trovati

#
        %ListaFile = ();                                                                                # hash contenente la lista dei file CDN
        opendir(my $dh, $localCDNDataDir) || die "Impossibile aprire $localCDNDataDir \n";
        # Scansione della directory per identificare i file da elaborare
        $FileCount = 0;
                while(readdir $dh) {
                #debug -->
                #print "File scansito: $_\t - \t";
                if( $_ =~ /^$myCDNFileTemplate.+\.gz$/) {
                        $FullPathFileName = $localCDNDataDir . '/' . $_;        # Costruisce il full path name
                        $ListaFile{$FullPathFileName} = 1;                                                                                      # Memorizza nell'hash il nome del file
                        $FileCount++;                                                                                           # Numero di file trovati
                        #debug -->
                        #print "File trovato: $ListaFile{$_} \n";
                }
                else {
                        #debug -->
                        #print "File non utilizzabile\n";
                        $UnusableFile++;
                }
        }
        closedir $dh;
        #debug -->
        print STDERR "\n **** Directory: $localCDNDataDir - File CDN trovati $FileCount - File Scartati $UnusableFile **** \n";
}