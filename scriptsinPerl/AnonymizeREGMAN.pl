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
use Switch;                                                                                                                                                     # Implementa lo statement Case - Switch
#++++++++++++++++++++++++++++
use constant STEP => 10000;                                                                                                                     # Passo di scrittura su riga di avanzamento
#++++++++++++++++++++++++++++
#use Date::Calc qw(Week_of_Year Day_of_Week Week_Number Day_of_Year Delta_DHMS);        # Libreria per manipolazione date
#use Carp;
#++++++++++++++++++++++++++++
# Dichiarazione subroutines
sub ScanCDNDir ($$);
#++++++++++++++++++++++++++++
# Legge il tipo di elaborazione richiesta
my $EvalutionType = shift or die "Usage: $0 tipo di elaborazione richiesta: AVS\|CDN\|REGMAN\n";
# Lista di parametri ammessi nella rivadi comando
%EvaluationList = ( 'AVS', '1', 'CDN', '0', 'REGMAN', '0' );
# Verifica il parametro
if ($EvaluationList{$EvalutionType} == 1){
        print " **** Inizio anonimizzazione file $EvalutionType ****\n"
} else
{
        die " **** Tipo di elaborazione non prevista: $EvalutionType ****\n"
}
# Dichiarazione directory e nomi file
# $myWorkdir = '/Volumes/gioMacBookL/Automation/Analytics/';                                            # Directory di lavoro - MAC
$myWorkdir = '/home/gioml/Analytics/';                                                                                          # Directory di lavoro - Server
# $myCuboAVSDataDir = '/Volumes/gioMacBookL/Automation/';                                                       # Directory dati origine  AVS e Trap - MAC
$myAVSDataDir = '/home/gioml/Analytics/AVS/';                                                                           # Directory dati origine  AVS - Server
$myTrapDataDir = '/home/gioml/Analytics/REGMAN/';                                                                       # Directory dati origine  TRAP da REGMAN - Server
# $myCDNDataDir[0]= '/Volumes/gioMacBookL/Automation/Analytics/h_20_UTC';                               # Directory dati di origine CDN ore 20 - MAC
# $myCDNDataDir[1]= '/Volumes/gioMacBookL/Automation/Analytics/h_21_UTC';                               # # Directory dati di origine CDN ore 21 -MAC
$myCDNDataDir[0]= '/home/gioml/Analytics/h_20_UTC';                                                                     # Directory dati di origine CDN ore 20 - Server
$myCDNDataDir[1]= '/home/gioml/Analytics/h_21_UTC';                                                                     # Directory dati di origine CDN ore 21 - Server
# File contenenti i dati di AVS e Trap
$myTrapFile = $myTrapDataDir . 'sample_cubo_traps_20181130.csv.zip';                            # File contenente le trap dei dispositivi (STB, Connected TV)
# $myAVSFile = $myAVSDataDir . 'Fruizioni_AVS_30112018.zip';                                                    # File contenente i dati di AVS (Front End di Tim Vision)
 $myAVSFile = $myAVSDataDir . 'prova_avs.csv.zip';
# Redirige stderr e impone autoflush
my $logfile = $myWorkdir. 'logerrAnonymize.txt';
open STDERR, ">$logfile" or die "\ **** Impossibile aprire file di errori $logfile *****\n";
local $| = 1;                                                                                                                                           # standard out autoflush
print "\n ++++++++ \nFile di log creato: $logfile ++++++++\n";
# Inizializzazione variabili per la visualizzazione dell'avanzamento attivitav =  STEP - 1;                                     # passo di avanzamento visualizzazione 1 milione di record
$index = 0;                                                     # per record elaborati == numero di linee lette
$myCounter = 0;                                         # per testo a capo
#++++++++++++++++++++++++++++
# Costruisce nome del file di uscita e lo apre - Dati CDN
# die ("Debug fase 1");
#
$CliStats = $myWorkdir.'AVS_TRAPS_20181130_anonymized.txt.gz';                                                                          # i file da elaborare contengono le trap dei STB e CTV
unlink $CliStats if -e $CliStats or die "\n ***** Impossibile cancellare il file $CliStats \n";         # creer file di uscita con le gli IP mascherati per i chunk erogati
my $gz_out = gzopen($CliStats, "wb") or die "\n **** Impossibile aprire $CliStats in scrittura: $! ****\n";
# Scrive l'header di uscita - Dati di AVS
$PrintedChars = 0;                                                                                                                                              # Numero di caratteri scritti nel file risultato
$Outputline ='Dispositivo|Timestamp|CLI|CPE-ID|ID_contenuto|Piattaforma|Vendor|Tipo_linea|Tempo_di_Visione';
$PrintedChars += $gz_out->gzwrite($Outputline);
#+++++++++++++++++++++++++++
# Time stamp - Valuta le prestazioni e visualizza avanzamento
($sec,$min,$hour,$mday,$mon,$year,$wday,$yday,$isdst)=localtime(time);
$StartTime = sprintf "%4d-%02d-%02d %02d:%02d:%02d",$year+1900,$mon+1,$mday,$hour,$min,$sec;
print STDERR "\n **** $StartTime - Google PoC - Anonimizzazione dati: Fonte AVS ****\n";
print "\n ++++++++ \n$StartTime> Avvio elaborazione Dati AVS\n";
print "\n                               - Milioni di record valutati - \n\n";
print "1       10        20       30        40        50        60        70        80       90       100\n";   #
#+++++++++++++++++++++++++++
# Ciclo su tutti i file delle directory contenenti i dati CDN
# Variabili di contesto
# $HeaderYesNO = 0;                                                                                                                                             # Variabile sentinella per la scrittura della riga di header
$myCrypt = Crypt::Lite->new( debug => 0 );                                                                                              # Instanzia la funzione di criptazione
$mySecret = 'tim_UFE1';                                                                                                                                 # Stringa con chiave di criptazione
# Inizia il ciclo
$ListaFile{$myAVSFile} = 1;                                                                                                                             # Memorizza nell'hash il nome del file contenente le fruizioni su AVS
foreach my $FileInput (sort keys %ListaFile)
{
        # Apre il file di ingresso, i-esimo elemento della lista
        my $gz_in = new IO::Uncompress::Unzip $FileInput [OPTS]
                or die "\n$FileInput unzip failed: Uncompress read error\n";
        # Ciclo sul file dati, riga per riga
#       if ($HeaderYesNO == 0) {$HeaderYesNO = 1; $PrintedChars += $gz_out->gzwrite($HeaderCDN);}       # Scrive l'intestazione una sola volta nel file risultato
        LOADLOOP: while ($line = $gz_in->getline())  {
                # print   "$line\n";                                                                                            # debug only
                @fields = split(/\t/, $line);                                                                                   # suddivisione della riga in stringhe separate dal carattere pipe
                $CLI = $fields[2];                                                                                              # CLI
                $CPE_ID = $fields[3];                                                                                                                   # CPE_ID
                # print "Record $ChunkCount - IP: $IPAddress\n";                                                                # debug only
                if ($CLI =~ /^[0-9]+/){                                                                                                 # Verifica la presenza dell'indirizzo IP
                        # print "\n **** Indirizzo IP: $fields[2] - Indirizzo Criptato: ";                      # debug only
                        $fields[2] = '0' . $fields[2];                                                                                          # Antepone lo '0' del prefisso
                        $fields[2] = $myCrypt->encrypt($fields[2], $mySecret);                                          # Cripta il CLI
                        # print "$fields[2] ****\n";                                                                                            # debug only
                        $fields[3] = $myCrypt->encrypt($fields[3], $mySecret);                                          # Cripta il CPE ID
                        $Outputline = "$fields[0]|$fields[1]|$fields[2]|$fields[3]|$fields[5]|$fields[12]|$fields[13]|$fields[16]|$fields[17]";  # Linea criptata, esclusi i campi non di interesse                                                         # Costruisce la riga per la scrittura nel file risultato
                        $PrintedChars += $gz_out->gzwrite($Outputline)
                                or die "\n **** Impossibile scrivere sul file $CliStats ****\n";                # Scrive la riga sul file di uscita, dopo l'anonimizzazione dell'indirizzo IP
                }
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
        close($gz_in) or die "\n **** Impossibile chiudere il file $myAVSFile ****\n";
        $ChunkCount++;                                                                                                                                  # Aggiorna il progressivo del file all'interno della settimana
        # $gz_out->gzclose();                                                                                                                   # Debug only
        # die "\n **** Esempio file cdn - $PrintedChars caratteri scritti ****\n";              # Debug only
}       # Fine ciclo sui nomi file
#
# Synthesis
print STDERR "\n **** AVS file Evaluated. $ChunkCount record found ****\n";
print "\nTerminata elaborazione dei file AVS - $ChunkCount record validi scritti su file risultato ****\n";
# Chiusura file di I/O
$gz_out->gzclose() or die "\n **** Impossibile chiudere il file $CliStats ****\n";
#
# Time Stamp
($sec,$min,$hour,$mday,$mon,$year,$wday,$yday,$isdst)=localtime(time);
$StopTime = sprintf "%4d-%02d-%02d %02d:%02d:%02d",$year+1900,$mon+1,$mday,$hour,$min,$sec;
print "\n **** Fine elaborazione alle $StopTime ****\n";
# Completa l'esecuzione
exit(1);
#
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