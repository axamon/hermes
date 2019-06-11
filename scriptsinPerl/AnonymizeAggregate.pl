#!/usr/bin/perl 
# Active State Perl   v5.12  for i386 su piattaforma Win XP/2000/2003 Server 
# Data:         06/12/2018
# Versione      1.0 Anonimizzazione dati per PoC Machine Learning con Google/Ingenia - Prima Versione
# Esempio riga da file sorgente
#             0                 1                        2                  3           4           5         6        7         8   9     10                11       12           13         14                     15   16          17                         
# CUBOTECHNICOLOR_ANDROIDTV     2018-11-30T00:00:08 000481808957 A7655449000686B9 CATALOGUE  50707453 STANDARD ABBONAMENTO |   |   |g.korenjak@alice.it|   |  CUBOVISION  |  ANDROID | Abbonamento;DownloadPlay |  | RETE FISSA | 48
# Esempio riga CDN
#           0                                                   1                                  2                                    3                               5                                   6                                                                                                                          8
# time-recvd(millisecond)       time-to-serve(microsecond)      client-ip       request-desc/response-status    bytes-sent      request-method  request-url        mine-type        req-header(user-agent)  abr-protocol    bitrate asset-id        session-id      entry-generated-time    client-type     profile
# [30/Nov/2018:19:59:59.786+0000]       218801                          79.18.78.113    TCP_MISS/200                             122508         GET                http://voddashhttps.cb.ticdn.it/videoteca2/V3/Film/2016/02/50541341/DASH/11475246/r16/191.mp4f   video/mp4       "Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko"     -       0       -       -
# Autore:       Giovanni Lepri
# Nome:                 Anonymize.pl
# Scopo:        Analisi del traffico dei clienti BB mobile su trial cap allo scopo di individuare profilature e comportamenti
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
#++++++++++++++++++++++++++++
use constant STEP => 10000;                                                                                                                     # Passo di scrittura su riga di avanzamento
#++++++++++++++++++++++++++++
#use Date::Calc qw(Week_of_Year Day_of_Week Week_Number Day_of_Year Delta_DHMS);        # Libreria per manipolazione date
#use Carp;
#++++++++++++++++++++++++++++
# Dichiarazione subroutines
sub ScanCDNDir ($$);
#++++++++++++++++++++++++++++
# Dichiarazione directory e nomi file
# $myWorkdir = '/Volumes/gioMacBookL/Automation/Analytics/';                                            # Directory di lavoro - MAC
$myWorkdir = '/home/gioml/Analytics/CDN/';                                                              # Directory di lavoro - Server
# $myCuboAVSDataDir = '/Volumes/gioMacBookL/Automation/';                                               # Directory dati origine  AVS e Trap - MAC
$myCuboAVSDataDir = '/home/gioml/Analytics/AVS/';                                                       # Directory dati origine  AVS e Trap - Server
# $myCDNDataDir[0]= '/Volumes/gioMacBookL/Automation/Analytics/h_20_UTC';                               # Directory dati di origine CDN ore 20 - MAC
# $myCDNDataDir[1]= '/Volumes/gioMacBookL/Automation/Analytics/h_21_UTC';                               # # Directory dati di origine CDN ore 21 -MAC
$myCDNDataDir[0]= '/home/gioml/Analytics/h_20_UTC';                                                     # Directory dati di origine CDN ore 20 - Server
$myCDNDataDir[1]= '/home/gioml/Analytics/h_21_UTC';                                                     # Directory dati di origine CDN ore 21 - Server
# File contenenti i dati di AVS e Trap
$myTrapFile = $myCuboAVSDataDir . 'sample_cubo_traps_20181130.csv.zip';                         # File contenente le trap dei dispositivi (STB, Connected TV)
$myAVSFile = $myCuboAVSDataDir . 'Fruizioni_30112018all.zip';                                           # File contenente i dati di AVS (Front End di Tim Vision)
# Redirige stderr e impone autoflush
my $logfile = $myWorkdir. 'logerrAnonymize.txt';
open STDERR, ">$logfile" or die "\ **** Impossibile aprire file di errori $logfile *****\n";
local $| = 1;                                                                                                                                           # standard out autoflush
print "\n ++++++++ \nFile di log creato: $logfile ++++++++\n";
#++++++++++++++++++++++++++++
# Costruisce la lista dei file contenente i dati CDN
$myCDNFileTemplate = 'we_accesslog_clf_';                                                                               # File contenente i log della CDN (Template)
ScanCDNDir ($myCDNDataDir[0], $myCDNFileTemplate);
ScanCDNDir ($myCDNDataDir[1], $myCDNFileTemplate);

#++++++++++++++++++++++++++++

# Inizializzazione variabili per la visualizzazione dell'avanzamento attivitav =  STEP - 1;                                     # passo di avanzamento visualizzazione 1 milione di record
$index = 0;                                                     # per record elaborati == numero di linee lette
$myCounter = 0;                                         # per testo a capo
#++++++++++++++++++++++++++++

# Time stamp - Valuta le prestazioni e visualizza avanzamento
($sec,$min,$hour,$mday,$mon,$year,$wday,$yday,$isdst)=localtime(time);
$StartTime = sprintf "%4d-%02d-%02d %02d:%02d:%02d",$year+1900,$mon+1,$mday,$hour,$min,$sec;
print STDERR "\n **** $StartTime - Google PoC - Anonimizzazione dati: Fonte CDN ****\n";
print "\n ++++++++ \n$StartTime> Avvio elaborazione Dati CDN\n";
print "\n                               - Milioni di record valutati - \n\n";
print "1       10        20       30        40        50        60        70        80       90       100\n";   #
#+++++++++++++++++++++++++++
# Ciclo su tutti i file delle directory contenenti i dati CDN
# Variabili di contesto
$HeaderYesNO = "Time;IP_Address;Content;Hit_Number;Miss_Number;Aborted_Number;AVG_Bit_Rate;Avg_Serving_Time;User_Agent";                # Variabile sentinella per la scrittura della riga di header
$PrintedChars = 0;                                                      #
$myCrypt = Crypt::Lite->new( debug => 0 );                                                                                              # Instanzia la funzione di criptazione
$mySecret = 'tim_UFE1';                                                                                                                                 # Stringa con chiave di criptazione
# Inizia il ciclo
foreach my $FileInput (sort keys %ListaFile)
{

        # Azzera gli hash con le statistiche
        %ip_hit_num = ();
        %ip_abort_num = ();
        %ip_miss_num = ();
        %ip_bytes = ();
        %ip_bitrate = ();
        %ip_time2serve = ();
        %ip_counter = ();
        %ip_user_agent = ();
        #+++++++++++++++++++++++++++
        # Apre il file di uscita - recupera il nome dall'hash ListaFileAggregati costruito in precedenza - Dati CDN
        # die ("Debug fase 1");
        #
        $CliStats = $ListaFileAggregati{$FileInput};
        # print "$CliStats \n ";
        unlink $CliStats if -e $CliStats; # or die "\n ***** Impossibile cancellare il file $CliStats \n";              # creer file di uscita con le gli IP mascherati per i chunk erogati
        my $gz_out = gzopen($CliStats, "wb") or die "\n **** Impossibile aprire $CliStats in scrittura: $! ****\n";
        $PrintedChars += $gz_out->gzwrite($HeaderYesNO);
        #+++++++++++++++++++++++++++
        # Apre il file di ingresso, i-esimo elemento della lista
        my $gz_in = gzopen($FileInput, "rb")                                                                                            # apre il file contenente i chunk cdn
                or die "\n$FileInput unzip failed: gzopen read error\n";
        # Ciclo sul file dati, riga per riga
        $gz_in->gzreadline($HeaderCDN);                                                                                                         # La prima riga contiene la release sw del logger
        $gz_in->gzreadline($HeaderCDN);                                                                                                         # La seconda riga contiene il nome dei campidel log
        # if ($HeaderYesNO == 0) {$HeaderYesNO = 1; $PrintedChars += $gz_out->gzwrite($HeaderCDN);}     # Scrive l'intestazione una sola volta nel file risultato
        LOADLOOP: while ($gz_in->gzreadline($line) > 0)  {
                # print   "$line\n";                                                                                            # debug only
                @fields = split(/\t/, $line);                                                                                   # suddivisione della riga in stringhe separate dal carattere pipe
                $IPAddress = $fields[2];                                                                                        # indirizzo IP
                if ($IPAddress =~ /^[0-9]{1,3}\.[0-9]{1,3}/){                                                                                   # Verifica la presenza e formato dell'indirizzo IP
                        @content_fields = split(/_/, $fields[6]);                                                                       # Identifica il contenuto
                        $Content = $content_fields[0];
                        # print "@fields; IP Address : $fields[2]; \n";
                        # print "Record $ChunkCount - IP: $IPAddress\n";                                                                # debug only
                        if ($fields[3] =~ /ABORTED/) {$ip_abort_num{$IPAddress}{$Content}++;}
                        elsif ($fields[3] =~ /MISS/){$ip_miss_num{$IPAddress}{$Content}++;}
                        else{$ip_hit_num{$IPAddress}{$Content}++;}
                        $ip_bytes{$IPAddress}{$Content} += $fields[4];
                        $ip_bitrate{$IPAddress}{$Content} += $fields[4]/$fields[1]*8;
                        $ip_counter{$IPAddress}{$Content}++ ;
                        $ip_user_agent{$IPAddress}{$Content} = $fields[8];
                        $ip_time2serve{$IPAddress}{$Content} += $fields[1];
                }
                # Visualizza l'avanzamento
                $index++;
                if ($index > $av) {
                        $av += STEP;
                        print STDOUT "+";                                                                                       # debug - Togliere il commento in caso di esecuzione
                        $myCounter++;
                        if ($myCounter > 98) {
                                $myCounter = 0;
                                print "\n";
                        }
                }
        }       # Fine esplorazione di un file
                # Calcolo bit_rate medio, tempo medio serving, tempo massimo e scrittura stringa
        foreach my $my_ip (sort keys %ip_bitrate) {
                foreach my $my_content (keys %{ $ip_bitrate{$my_ip} }) {
                        # debug only
                        #print "$my_ip, $my_content: $ip_bitrate{$my_ip}{$my_content} $ip_counter{$my_ip}{$my_content}; \n";
                        $ip_bitrate{$my_ip}{$my_content} = $ip_bitrate{$my_ip}{$my_content}/$ip_counter{$my_ip}{$my_content};
                }
        }
        foreach my $my_ip (sort keys %ip_time2serve) {
                foreach my $my_content (keys %{ $ip_time2serve{$my_ip} }) {
                        $ip_time2serve{$my_ip}{$my_content} = $ip_time2serve{$my_ip}{$my_content}/$ip_counter{$my_ip}{$my_content};
                        # Criptazione
                        $my_ip_enc = $myCrypt->encrypt($my_ip, $mySecret);                                              # Cripta l'indirizzo IP
                        # "Time;IP_Address;Content;Hit_Number;Miss_Number;Aborted_Number;AVG_Bit_Rate;Avg_Serving_Time;User_Agent";
$OutputLine = "$Ora{$FileInput};$my_ip_enc;$my_content;$ip_hit_num{$my_ip}{$my_content};$ip_miss_num{$my_ip}{$my_content};$ip_abort_num{$my_ip}{$my_content};$ip_bitrate{$my_ip}{$my_content};$ip_time2serve{$my_ip}{$my_content};$ip_user_agent{$my_ip}{$my_content}\n";
                        # debug only
                        # print "$OutputLine \n";
                        $PrintedChars += $gz_out->gzwrite($OutputLine)
                                or die "\n **** Impossibile scrivere sul file $CliStats ****\n";
                }
        }

        # Chiude il file di I/O
        $gz_in->gzclose();
        $ChunkCount++;                                                                          # Aggiorna il progressivo del file all'interno della settimana
        # $gz_out->gzclose();                                                                   # Debug only
        # die "\n **** Esempio file cdn - $PrintedChars caratteri scritti ****\n";              # Debug only
}       # Fine ciclo sui nomi file
#
# Synthesis
print STDERR "\n **** CDN files Evaluated. $ChunkCount chunks found ****\n";
print "\nTerminata elaborazione dei file CDN - $ChunkCount chunk validi scritti su file risultato ****\n";
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

        my @filenamefields;
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
                        $ListaFile{$FullPathFileName} = 1;                      # Memorizza nell'hash il nome del file
                        $_ =~ /(^$myCDNFileTemplate.+)\.gz$/;
                        $ListaFileAggregati{$FullPathFileName} = $myWorkdir . $1 . '_aggr.gz';
                        @filenamefields = split(/_/, $FullPathFileName);
                        $Ora{$FullPathFileName} = $filenamefields[7];
                        $FileCount++;                                           # Numero di file trovati
                        #debug -->
                        #print "File trovato: $ListaFile{$_} \n";
                        #print "$ListaFileAggregati{$FullPathFileName}, $Ora{$FullPathFileName} \n";
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